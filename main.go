package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/mpetavy/common"
	"regexp"
	"slices"
	"sort"
	"strings"
)

//go:embed go.mod
var resources embed.FS

type OM[V any] struct {
	table map[string]V
	keys  []string
}

func NewOM[V any]() OM[V] {
	return OM[V]{
		table: make(map[string]V, 0),
		keys:  make([]string, 0),
	}
}

func (o *OM[V]) Set(k string, v V) {
	o.table[k] = v
	o.keys = append(o.keys, k)
}

func (o OM[V]) MarshalJSON() ([]byte, error) {
	ba, err := json.Marshal(o.table)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile("\".*?[,}]")

	buf := bytes.Buffer{}
	buf.WriteString("{")

	elements := regex.FindAllString(string(ba), -1)

	if len(elements) == 0 {
		buf.WriteString("}")
	} else {
		sort.Slice(elements, func(i, j int) bool {
			k0 := elements[i][1:]
			k0 = k0[:strings.Index(k0, "\"")]
			k1 := elements[j][1:]
			k1 = k1[:strings.Index(k1, "\"")]

			return slices.Index(o.keys, k0) < slices.Index(o.keys, k1)
		})

		for i := 0; i < len(elements); i++ {
			buf.WriteString(elements[i][0 : len(elements[i])-1])

			if i+1 == len(elements) {
				buf.WriteString("}")
			} else {
				buf.WriteString(",")
			}
		}
	}

	return buf.Bytes(), nil
}

func (o *OM[V]) UnmarshalJSON(data []byte) error {
	clear(o.table)
	clear(o.keys)

	err := json.Unmarshal(data, &o.table)
	if err != nil {
		return err
	}

	regex := regexp.MustCompile("[\\{,]\".*?\"")

	keys := regex.FindAllString(string(data), -1)
	for _, key := range keys {
		key = key[2 : len(key)-1]

		o.keys = append(o.keys, key)
	}

	return nil
}

func init() {
	common.Init("", "", "", "", "", "", "", "", &resources, nil, nil, run, 0)
}

func run() error {
	om := NewOM[string]()

	om.Set("c", "cc")
	om.Set("b", "bb")
	om.Set("a", "aa")

	//om := NewOM[int]()
	//
	//om.Set("c", 3)
	//om.Set("b", 2)
	//om.Set("a", 1)

	for _, k := range om.keys {
		fmt.Printf("%v: %v\n", k, om.table[k])
	}

	ba, err := json.Marshal(om)
	if common.Error(err) {
		return err
	}

	fmt.Printf("%s\n", ba)

	om2 := NewOM[string]()
	err = json.Unmarshal(ba, &om2)
	if common.Error(err) {
		return err
	}

	for _, k := range om2.keys {
		fmt.Printf("%v: %v\n", k, om2.table[k])
	}

	return nil
}

func main() {
	common.Run(nil)
}
