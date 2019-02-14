package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/levigross/grequests"
	"github.com/malice-plugins/pkgs/database"
	"github.com/malice-plugins/pkgs/database/elasticsearch"
	"github.com/malice-plugins/pkgs/utils"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	name     = "comodo"
	category = "av"
)

var (
	// Version stores the plugin's version
	Version string
	// BuildTime stores the plugin's build time
	BuildTime string

	path string

	// es is the elasticsearch database object
	es elasticsearch.Database
)

type pluginResults struct {
	ID   string      `json:"id" structs:"id,omitempty"`
	Data ResultsData `json:"comodo" structs:"comodo"`
}

// Comodo json object
type Comodo struct {
	Results ResultsData `json:"comodo"`
}

// ResultsData json object
type ResultsData struct {
	Infected bool   `json:"infected" structs:"infected"`
	Result   string `json:"result" structs:"result"`
	Engine   string `json:"engine" structs:"engine"`
	Updated  string `json:"updated" structs:"updated"`
	MarkDown string `json:"markdown,omitempty" structs:"markdown,omitempty"`
}

func assert(err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"plugin":   name,
			"category": category,
			"path":     path,
		}).Fatal(err)
	}
}

// AvScan performs antivirus scan
func AvScan(path string, timeout int) Comodo {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	output, err := utils.RunCommand(ctx, "/opt/COMODO/cmdscan", "-vs", path)
	assert(err)

	return Comodo{
		Results: ParseComodoOutput(output),
	}
}

// ParseComodoOutput convert comodo output into ResultsData struct
func ParseComodoOutput(comodoout string) ResultsData {

	comodo := ResultsData{Infected: false, Engine: getComodoVersion(), Updated: getUpdatedDate()}

	log.Debug("comodoout: ", comodoout)

	// EXAMPLE OUTPUT:
	// -----== Scan Start ==-----
	// /malware/EICAR ---> Found Virus, Malware Name is Malware
	// -----== Scan End ==-----
	// Number of Scanned Files: 1
	// Number of Found Viruses: 1
	lines := strings.Split(comodoout, "\n")

	// Extract Virus string
	if len(lines[1]) != 0 {
		if strings.Contains(lines[1], "Found Virus") {
			result := extractVirusName(lines[1])
			if len(result) != 0 {
				comodo.Result = result
				comodo.Infected = true
				return comodo
			}
			fmt.Println("[ERROR] Virus name extracted was empty: ", result)
			os.Exit(2)
		}
	}

	return comodo
}

// extractVirusName extracts Virus name from scan results string
func extractVirusName(line string) string {
	keyvalue := strings.Split(line, "is")
	return strings.TrimSpace(keyvalue[1])
}

func updateAV() error {
	fmt.Println("Updating Comodo...")
	response, err := grequests.Get("http://download.comodo.com/av/updates58/sigs/bases/bases.cav", nil)
	if err != nil {
		return err
	}

	if response.Ok != true {
		log.Println("Request did not return OK")
	}

	if err = response.DownloadToFile("/opt/COMODO/scanners/bases.cav"); err != nil {
		log.Println("Unable to download file: ", err)
	}
	// Update UPDATED file
	t := time.Now().Format("20060102")
	err = ioutil.WriteFile("/opt/malice/UPDATED", []byte(t), 0644)
	return err
}

func getComodoVersion() string {
	file, err := os.Open("/opt/COMODO/etc/COMODO.xml")
	assert(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "<ProductVersion>") {
			versionOut := strings.TrimSpace(strings.Replace(strings.Replace(line, "<ProductVersion>", "", 1), "</ProductVersion>", "", 1))
			log.Debug("Comodo Version: ", versionOut)
			return versionOut
		}
	}
	return "error"
}

func getUpdatedDate() string {
	if _, err := os.Stat("/opt/malice/UPDATED"); os.IsNotExist(err) {
		return BuildTime
	}
	updated, err := ioutil.ReadFile("/opt/malice/UPDATED")
	assert(err)
	return string(updated)
}

func parseUpdatedDate(date string) string {
	layout := "200601021504"
	t, _ := time.Parse(layout, date)
	return fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
}

func generateMarkDownTable(c Comodo) string {
	var tplOut bytes.Buffer

	t := template.Must(template.New("comodo").Parse(tpl))

	err := t.Execute(&tplOut, c)
	if err != nil {
		log.Println("executing template:", err)
	}

	return tplOut.String()
}

func printStatus(resp gorequest.Response, body string, errs []error) {
	fmt.Println(body)
}

func webService() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/scan", webAvScan).Methods("POST")
	log.Info("web service listening on port :3993")
	log.Fatal(http.ListenAndServe(":3993", router))
}

func webAvScan(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(32 << 20)
	file, header, err := r.FormFile("malware")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Please supply a valid file to scan.")
		log.Error(err)
	}
	defer file.Close()

	log.Debug("Uploaded fileName: ", header.Filename)

	tmpfile, err := ioutil.TempFile("/malware", "web_")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	data, err := ioutil.ReadAll(file)
	assert(err)

	if _, err = tmpfile.Write(data); err != nil {
		log.Fatal(err)
	}
	if err = tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	// Do AV scan
	comodo := AvScan(tmpfile.Name(), 60)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(comodo); err != nil {
		log.Fatal(err)
	}
}

func main() {

	cli.AppHelpTemplate = utils.AppHelpTemplate
	app := cli.NewApp()

	app.Name = "comodo"
	app.Author = "blacktop"
	app.Email = "https://github.com/blacktop"
	app.Version = Version + ", BuildTime: " + BuildTime
	app.Compiled, _ = time.Parse("20060102", BuildTime)
	app.Usage = "Malice AVG AntiVirus Plugin"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "verbose output",
		},
		cli.StringFlag{
			Name:        "elasticsearch",
			Value:       "",
			Usage:       "elasticsearch url for Malice to store results",
			EnvVar:      "MALICE_ELASTICSEARCH_URL",
			Destination: &es.URL,
		},
		cli.BoolFlag{
			Name:  "table, t",
			Usage: "output as Markdown table",
		},
		cli.BoolFlag{
			Name:   "callback, c",
			Usage:  "POST results back to Malice webhook",
			EnvVar: "MALICE_ENDPOINT",
		},
		cli.BoolFlag{
			Name:   "proxy, x",
			Usage:  "proxy settings for Malice webhook endpoint",
			EnvVar: "MALICE_PROXY",
		},
		cli.IntFlag{
			Name:   "timeout",
			Value:  60,
			Usage:  "malice plugin timeout (in seconds)",
			EnvVar: "MALICE_TIMEOUT",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "update",
			Usage: "Update virus definitions",
			Action: func(c *cli.Context) error {
				return updateAV()
			},
		},
		{
			Name:  "web",
			Usage: "Create a Comodo scan web service",
			Action: func(c *cli.Context) error {
				webService()
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {

		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}

		if c.Args().Present() {
			path, err := filepath.Abs(c.Args().First())
			assert(err)

			if _, err := os.Stat(path); os.IsNotExist(err) {
				assert(err)
			}

			comodo := AvScan(path, c.Int("timeout"))
			comodo.Results.MarkDown = generateMarkDownTable(comodo)

			// upsert into Database
			if len(c.String("elasticsearch")) > 0 {
				err := es.Init()
				if err != nil {
					return errors.Wrap(err, "failed to initalize elasticsearch")
				}
				err = es.StorePluginResults(database.PluginResults{
					ID:       utils.Getopt("MALICE_SCANID", utils.GetSHA256(path)),
					Name:     name,
					Category: category,
					Data:     structs.Map(comodo.Results),
				})
				if err != nil {
					return errors.Wrapf(err, "failed to index malice/%s results", name)
				}
			}

			if c.Bool("table") {
				fmt.Println(comodo.Results.MarkDown)
			} else {
				comodo.Results.MarkDown = ""
				avgJSON, err := json.Marshal(comodo)
				assert(err)
				if c.Bool("callback") {
					request := gorequest.New()
					if c.Bool("proxy") {
						request = gorequest.New().Proxy(os.Getenv("MALICE_PROXY"))
					}
					request.Post(os.Getenv("MALICE_ENDPOINT")).
						Set("X-Malice-ID", utils.Getopt("MALICE_SCANID", utils.GetSHA256(path))).
						Send(string(avgJSON)).
						End(printStatus)

					return nil
				}
				fmt.Println(string(avgJSON))
			}
		} else {
			log.Fatal(fmt.Errorf("Please supply a file to scan with malice/comodo"))
		}
		return nil
	}

	err := app.Run(os.Args)
	assert(err)
}
