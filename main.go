package main

import (
	"embed"
	"github.com/mpetavy/common"
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	return nil
}

func main() {
	common.Run(nil)
}
