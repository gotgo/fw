package multi

import cphash "github.com/kavu/go-phash"

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
	hash, err := cphash.ImageHashDCT(filepath)
	return hash, err
}
