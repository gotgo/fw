package logging

import (
	"errors"
	"io"

	"github.com/cihub/seelog"
)

type KV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func SetKeyValue(lm *LogMessage, kv ...*KV) {
	if kv == nil || len(kv) == 0 {
		return
	} else if len(kv) == 1 {
		lm.Key = kv[0].Key
		lm.Value = kv[0].Value
	} else {
		lm.Key = "data"
		lm.Value = kv
	}
}

func init() {
	DisableLog()
}

// DisableLog disables all library log output
func DisableLog() {
	logger = seelog.Disabled
}

// UseLogger uses a specified seelog.LoggerInterface to output library log.
// Use this func if you are using Seelog logging system in your app.
func UseLogger(newLogger seelog.LoggerInterface) {
	logger = newLogger
	newLogger.SetAdditionalStackDepth(2)
}

// SetLogWriter uses a specified io.Writer to output library log.
// Use this func if you are not using Seelog logging system in your app.
func SetLogWriter(writer io.Writer) error {
	if writer == nil {
		return errors.New("Nil writer")
	}

	newLogger, err := seelog.LoggerFromWriterWithMinLevel(writer, seelog.TraceLvl)
	if err != nil {
		return err
	}

	UseLogger(newLogger)
	return nil
}

// Call this before app shutdown
func FlushLog() {
	logger.Flush()
}

type RequstLogger interface {
	Capture(service string, name string, args map[string]string, duration int, outcome bool)
}

type Kind string

const (
	Error     = "error"
	Warn      = "warn"
	Inform    = "inform"
	Debug     = "debug"
	Timeout   = "timeout"
	Connect   = "connect"
	Event     = "event"
	Marshal   = "marshal"
	Unmarshal = "unmarshal"
	Panic     = "panic"
)

func NewLogger() Logger {
	return new(StdLogger)
}
