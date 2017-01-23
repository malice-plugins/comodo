package main_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	log "github.com/Sirupsen/logrus"
)

const resultString = `-----== Scan Start ==-----
/malware/EICAR ---> Found Virus, Malware Name is Malware
-----== Scan End ==-----
Number of Scanned Files: 1
Number of Found Viruses: 1
`

const versionString = `
`

func parseSophosVersion(versionOut string) (version string, database string) {

	lines := strings.Split(versionOut, "\n")

	for _, line := range lines {

		if strings.Contains(line, "Product version") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				version = strings.TrimSpace(parts[1])
			} else {
				log.Error("Umm... ", parts)
			}
		}

		if strings.Contains(line, "Virus data version") {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				database = strings.TrimSpace(parts[1])
				break
			} else {
				log.Error("Umm... ", parts)
			}
		}
	}

	return
}

func parseComodoOutput(comodoout string) (string, error) {

	comodo := ResultsData{Infected: false, Engine: "1.1"}
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

	comodo.Updated = getUpdatedDate()

	return comodo, nil
}

// TestParseResult tests the ParseFSecureOutput function.
func TestParseResult(t *testing.T) {

	results, err := parseComodoOutput(resultString)

	if err != nil {
		t.Log(err)
	}

	if true {
		t.Log("results: ", results)
	}

}

// TestParseVersion tests the GetFSecureVersion function.
func TestParseVersion(t *testing.T) {

	version, database := parseSophosVersion(versionString)

	if true {
		t.Log("version: ", version)
		t.Log("database: ", database)
	}

}
