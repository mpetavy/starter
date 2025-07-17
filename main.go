package main

import (
	"embed"
	"fmt"
	"github.com/mpetavy/common"
	"sync"
	"sync/atomic"
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

type Sperr struct {
	Current atomic.Int64

	mu sync.Mutex
}

var sp Sperr

func (sp *Sperr) CanEnter() bool {
	id := int64(common.GoRoutineId())

	if sp.Current.CompareAndSwap(id, id) {
		return false
	}

	sp.mu.Lock()

	sp.Current.Store(id)

	return true
}

func Error(msg error) bool {
	if !sp.CanEnter() {
		return msg != nil
	}

	defer sp.mu.Unlock()

	Error(fmt.Errorf("DONT SHOW ERROR!!"))

	fmt.Printf("%s\n", msg)

	return msg != nil
}

func run() error {
	Error(fmt.Errorf("test"))

	return nil
}

func main() {
	common.Run(nil)
}
