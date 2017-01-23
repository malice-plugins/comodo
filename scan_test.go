package main_test

import (
	"fmt"
	"strings"
	"testing"
)

const resultString = `-----== Scan Start ==-----
/malware/EICAR ---> Found Virus, Malware Name is Malware
-----== Scan End ==-----
Number of Scanned Files: 1
Number of Found Viruses: 1
`

func extractVirusName(line string) string {
	keyvalue := strings.Split(line, "is")
	return strings.TrimSpace(keyvalue[1])
}

func parseComodoOutput(comodoout string) (string, error) {

	lines := strings.Split(comodoout, "\n")

	// Extract Virus string
	if len(lines[1]) != 0 {
		if strings.Contains(lines[1], "Found Virus") {
			result := extractVirusName(lines[1])
			if len(result) != 0 {
				return result, nil
			} else {
				return "", fmt.Errorf("emptry string")
			}
		}
	}
	return "", fmt.Errorf("emptry string")
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
