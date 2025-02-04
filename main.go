package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"path/filepath"
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	jobFiles, err := common.ListFiles("/home/ransom/go/src/hakodate/testdata/events/*-job.json", false)
	if common.Error(err) {
		return err
	}

	for _, jobFile := range jobFiles {
		fmt.Printf("%s\n", jobFile)

		ba, err := os.ReadFile(jobFile)
		if common.Error(err) {
			return err
		}

		m := make(map[string]any)
		err = json.Unmarshal(ba, &m)
		if common.Error(err) {
			return err
		}

		sopInstanceUID, _ := m["sopInstanceUID"].(string)

		fmt.Printf("%s\n", sopInstanceUID)

		prefix := jobFile[:len(jobFile)-len("-job.json")]

		eventFile := prefix + "-event.json"
		responseFile := prefix + "-response.json"

		newJobFile := filepath.Join(filepath.Dir(jobFile), sopInstanceUID+"-job.json")
		newEventFile := filepath.Join(filepath.Dir(jobFile), sopInstanceUID+"-event.json")
		newResponseFile := filepath.Join(filepath.Dir(jobFile), sopInstanceUID+"-response.json")

		err = os.Rename(jobFile, newJobFile)
		if common.Error(err) {
			return err
		}

		err = os.Rename(eventFile, newEventFile)
		if common.Error(err) {
			return err
		}

		err = os.Rename(responseFile, newResponseFile)
		if common.Error(err) {
			return err
		}
	}

	return nil
}

func main() {
	common.Run(nil)
}
