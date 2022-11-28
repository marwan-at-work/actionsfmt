// Package actionsfmt helps run go commands that are pretty-formatted
// for GitHub actions. Currently it expects "go test" out of the box
// and formats it according to the GitHub Actions UI workflow commands:
// https://docs.github.com/en/actions/reference/workflow-commands-for-github-actions
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/logrusorgru/aurora/v3"
)

type testEvent struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

type testRun struct {
	name        string
	lines       []string
	finalAction string
}

func main() {
	var (
		dec       = json.NewDecoder(os.Stdin)
		mp        = map[string]*testRun{}
		testOrder []string
		failed    bool
	)

Outer:
	for dec.More() {
		var te testEvent
		err := dec.Decode(&te)
		must(err)
		if te.Test == "" {
			continue
		}
		if te.Action == "run" {
			testOrder = append(testOrder, te.Test)
			mp[te.Test] = &testRun{name: te.Test}
			continue
		}
		tr := mp[te.Test]
		switch te.Action {
		case "output":
			trimmed := strings.TrimSpace(te.Output)
			if strings.HasPrefix(trimmed, "===") || strings.HasPrefix(trimmed, "---") {
				// === and --- prefixed outputs are redundant information that we already collected
				continue Outer
			}
			lastChar := byte('\n')
			if len(tr.lines) >= 1 {
				lastLine := tr.lines[len(tr.lines)-1]
				lastChar = lastLine[len(lastLine)-1]
			}
			if lastChar != '\n' {
				tr.lines[len(tr.lines)-1] += te.Output
			} else {
				tr.lines = append(tr.lines, te.Output)
			}
		case "pass":
			tr.finalAction = "pass"
		case "fail":
			failed = true
			tr.finalAction = "fail"
		}
	}
	for _, testName := range testOrder {
		tr := mp[testName]
		sprint := fmt.Sprintf
		if tr.finalAction == "pass" {
			sprint = func(format string, a ...interface{}) string {
				return aurora.Green(fmt.Sprintf(format, a...)).String()
			}
		} else if tr.finalAction == "fail" {
			sprint = func(format string, a ...interface{}) string {
				return aurora.Red(fmt.Sprintf(format, a...)).String()
			}
		}
		fmt.Printf("::group::%s\n", sprint(tr.name))
		if len(tr.lines) == 0 {
			fmt.Println("All tests have passed")
		}
		for _, line := range tr.lines {
			fmt.Print(line)
		}
		fmt.Printf("\n::endgroup::\n")
	}
	if failed {
		log.Fatal(aurora.Red("Go Test Has Failed"))
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
