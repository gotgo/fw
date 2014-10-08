package io

type NopWriter struct{}

func (now *NopWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
