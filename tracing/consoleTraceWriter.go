package tracing

import "os"

type ConsoleTraceWriter struct {
}

func (ctw *ConsoleTraceWriter) Write(b []byte) (n int, err error) {
	return os.Stdout.Write(b)
}

func (ctw *ConsoleTraceWriter) Close() error {
	return nil
}
