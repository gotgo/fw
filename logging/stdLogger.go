package logging

import (
	"encoding/json"
	"fmt"

	seelog "github.com/cihub/seelog"
)

//even though Logger is an instance, all instances are using the same
//singleton logger. it's unclear if this should stay this way
var logger seelog.LoggerInterface

type StdLogger struct {
}

func (l *StdLogger) Log(m *LogMessage) {
	bytes, err := json.Marshal(m)
	if err != nil {
		l.MarshalFail("could not log event because of marshal fail", m, err)
		return
	}
	str := string(bytes)
	switch m.Kind {
	case Inform, Event:
		logger.Info(str)
	case Debug:
		logger.Debug(str)
	case Warn, Timeout:
		logger.Warn(str)
	case Error, Marshal, Unmarshal, Connect:
		logger.Error(str)
	case Panic:
		logger.Critical(str)
	default:
		logger.Info(str)
	}
}

func (l *StdLogger) MarshalFail(m string, obj interface{}, err error) {
	msg := err.Error()
	lm := &LogMessage{
		Message: "Marshal Failed " + m,
		Error:   msg,
		Key:     "object",
		Value:   fmt.Sprintf("%#v", obj),
		Kind:    "marshalFail",
	}
	bytes, err := json.Marshal(lm)
	if err != nil {
		l.MarshalFail("could not log event because of marshal fail", m, err)
		return
	}
	logger.Error(string(bytes))
}

func (l *StdLogger) UnmarshalFail(m string, data []byte, err error) {
	var persistedData []byte
	const arbitraryCutoffSize = 5000
	if len(data) < arbitraryCutoffSize {
		persistedData = data
	}

	msg := err.Error()
	lm := &LogMessage{
		Message: "Unmarshal Failed" + m,
		Error:   msg,
		Key:     "data",
		Value:   persistedData,
		Kind:    "unmarshalFail",
	}
	bytes, err := json.Marshal(lm)
	if err != nil {
		l.MarshalFail("could not log event because of marshal fail", m, err)
		return
	}
	logger.Error(string(bytes))
}

func (l *StdLogger) Timeout(m string, err error, kv ...*KeyValue) {
	msg := err.Error()
	lm := &LogMessage{
		Message: m,
		Error:   msg,
		Kind:    "timeout",
	}
	SetKeyValue(lm, kv...)
	bytes, err := json.Marshal(lm)
	if err != nil {
		l.MarshalFail("could not log event because of marshal fail", m, err)
		return
	}
	logger.Warn(string(bytes))
}

func (l *StdLogger) ConnectFail(m string, err error, kv ...*KeyValue) {
	msg := err.Error()
	lm := &LogMessage{
		Message: m,
		Error:   msg,
		Kind:    "connectFail",
	}
	SetKeyValue(lm, kv...)
	bytes, err := json.Marshal(lm)
	if err != nil {
		l.MarshalFail("could not log event because of marshal fail", m, err)
		return
	}
	logger.Warn(string(bytes))
}

func (l *StdLogger) WillPanic(m string, err error, kv ...*KeyValue) {
	msg := err.Error()
	lm := &LogMessage{
		Message: m,
		Error:   msg,
		Kind:    "willPanic",
	}
	SetKeyValue(lm, kv...)
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not marshal will panic LogMessage", lm, err)
	} else {
		logger.Critical(string(bytes))
	}
	logger.Flush()
}

func (l *StdLogger) HadPanic(m string, r interface{}) {
	//figure out what r is
	err, _ := r.(error)
	str, _ := r.(string)

	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	lm := &LogMessage{
		Message: m,
		Error:   errMsg,
		Value:   str,
		Kind:    "panic",
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not marshal panic LogMessage", lm, err)
	} else {
		logger.Critical(string(bytes))
	}
	logger.Flush()
}

func (l *StdLogger) Error(m string, e error, kv ...*KeyValue) {
	msg := "unknown error"
	if e != nil {
		msg = e.Error()
	}
	lm := &LogMessage{
		Message: m,
		Error:   msg,
		Kind:    "error",
	}
	SetKeyValue(lm, kv...)
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log error because of marshal fail, from error "+m, lm, err)
		return
	} else {
		logger.Error(string(bytes))
	}
}

func (l *StdLogger) Warn(m string, kv ...*KeyValue) {

	lm := &LogMessage{
		Message: m,
		Kind:    "warn",
	}
	SetKeyValue(lm, kv...)

	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log warn because of marshal fail", lm, err)
		return
	} else {
		logger.Warn(string(bytes))
	}
}

// Infom captures a simple message. If you are logging key value pairs,
// use Info(m interface{})
func (l *StdLogger) Inform(m string) {
	lm := &LogMessage{Message: m, Kind: "inform"}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("Could not log event because info message marshal fail", lm, err)
	} else {
		logger.Info(string(bytes))
	}
}

// Info logs key value pairs, typically to JSON. Typically using an anonymous struct:
//
//		log.Info(struct{MyKey string}{MyKey:"value to capture"})
func (l *StdLogger) Event(m string, kv ...*KeyValue) {
	lm := &LogMessage{
		Message: m,
		Kind:    "event",
	}
	SetKeyValue(lm, kv...)
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log event because of marshal fail", lm, err)
		return
	} else {
		logger.Info(string(bytes))
	}
}

func (l *StdLogger) Debug(m string, kv ...*KeyValue) {
	lm := &LogMessage{
		Message: m,
		Kind:    "debug",
	}

	SetKeyValue(lm, kv...)
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log event because of marshal fail", lm, err)
		return
	} else {
		logger.Debug(string(bytes))
	}
}
