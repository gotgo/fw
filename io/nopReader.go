package io

type NopReader struct{}

func (nor *NopReader) Read(p []byte) (n int, err error) {
	return 0, nil
}
