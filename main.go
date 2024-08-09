package main

import (
	"context"
	"embed"
	"flag"
	"github.com/mpetavy/common"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

//go:embed go.mod
var resources embed.FS

var (
	dryRun     = flag.Bool("n", true, "dry run")
	oldPath    = flag.String("old", "/home/ransom/temp1", "old path to be modified")
	targetPath = flag.String("new", "/home/ransom/go/src/hakodate", "new target path as template")
)

type Files = orderedmap.OrderedMap[string, string]

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

func walkFile(root string) (*Files, error) {
	files := orderedmap.New[string, string]()

	err := common.WalkFiles(root, true, false, func(path string, f os.FileInfo) error {
		if f.IsDir() && strings.HasPrefix(f.Name(), ".") {
			return fs.SkipDir
		}

		if path == root {
			return nil
		}

		files.Set(f.Name(), path)

		return nil
	})
	if common.Error(err) {
		return nil, err
	}

	return files, nil
}

func move(from, to string) error {
	common.Info("%s -> %s\n", from, to)

	targetDir := filepath.Join(*oldPath, filepath.Dir(to))
	if targetDir != "." && !common.FileExists(targetDir) && !*dryRun {
		err := os.MkdirAll(targetDir, common.DefaultDirMode)
		if common.Error(err) {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	cmds := []string{"mv", from, to}
	if *dryRun {
		cmds = slices.Insert(cmds, 1, "-n")
	}

	cmd := exec.CommandContext(ctx, "git", cmds...)
	cmd.Dir = *oldPath

	common.Debug(common.CmdToString(cmd))

	ba, err := cmd.CombinedOutput()
	if ba != nil {
		common.Debug("%s", ba)
	}
	if common.Error(err) {
		return err
	}

	return nil
}

func run() error {
	originalFiles, err := walkFile(*oldPath)
	if common.Error(err) {
		return err
	}
	targetFiles, err := walkFile(*targetPath)
	if common.Error(err) {
		return err
	}

	for pair := originalFiles.Oldest(); pair != nil; pair = pair.Next() {
		filename := pair.Key
		path := pair.Value

		found, ok := targetFiles.Get(filename)
		if !ok {
			continue
		}

		if path[len(*oldPath):] != found[len(*targetPath):] {
			target := found[len(*targetPath)+1:]

			err := move(path[len(*oldPath)+1:], target)
			if common.Error(err) {
				return err
			}
		}
	}

	return nil
}

func main() {
	common.Run(nil)
}
