package multi

import (
	"os"
	"time"

	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
	cphash "github.com/kavu/go-phash"
)

func PHashTaskIn(filepath string) interface{} {
	return filepath
}

func PHashTaskOut(out interface{}) uint64 {
	return out.(uint64)
}

type PHashTask struct {
	Log logging.Logger
}

func (p *PHashTask) Run(input interface{}) (interface{}, error) {
	filepath, ok := input.(string)
	if !ok {
		panic("wrong type")
	}
	time.Sleep(time.Second * 5)
	fi, err := os.Stat(filepath)
	if err != nil {
		me.LogError(p.Log, "failed with path "+filepath, err)
	}
	if fi.Size() == int64(0) {
		me.LogInform(p.Log, "size is zero sleeping.....for 10 seconds for "+filepath)
		time.Sleep(time.Second * 10)
	}

	if fi.Size() == int64(0) {
		me.LogInform(p.Log, "still size zero")
		return nil, me.NewErr("file is size zero")
	}

	hash, err := cphash.ImageHashDCT(filepath)
	return hash, err
}

func (p *PHashTask) Name() string {
	return "phash"
}
