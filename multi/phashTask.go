package multi

import (
	"fmt"
	"os"
	"time"

	cphash "github.com/kavu/go-phash"
)

func PHashTaskIn(filepath string) interface{} {
	return filepath
}

func PHashTaskOut(out interface{}) uint64 {
	return out.(uint64)
}

type PHashTask struct{}

func (p *PHashTask) Run(input interface{}) (interface{}, error) {
	filepath, ok := input.(string)
	if !ok {
		panic("wrong type")
	}
	fi, err := os.Stat(filepath)
	if err != nil || fi.Size() == int64(0) {
		fmt.Println("sleeping.....")
		time.Sleep(time.Second * 2)
	}

	hash, err := cphash.ImageHashDCT(filepath)
	return hash, err
}

func (p *PHashTask) Name() string {
	return "phash"
}
