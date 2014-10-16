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
	switch m.Kind {
	case Inform, Event:
		logger.Info(m)
	case Debug:
		logger.Debug(m)
	case Warn, Timeout:
		logger.Warn(m)
	case Error, Marshal, Unmarshal, Connect:
		logger.Error(m)
	case Panic:
		logger.Critical(m)
	default:
		logger.Info(m)
	}
}

func (l *StdLogger) MarshalFail(m string, obj interface{}, err error) {
	msg := err.Error()
	lm := &LogMessage{
		Message: "Marshal Failed " + m,
		Error:   msg,
		Key:     "object",
		Value:   fmt.Sprintf("%#v", obj),
	}
	logger.Error(lm)
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
	}
	logger.Error(lm)
}

func (l *StdLogger) Timeout(m string, err error) {
	msg := err.Error()
	lm := &LogMessage{
		Message: "Timeout: " + m,
		Error:   msg,
	}
	logger.Warn(lm)
}

func (l *StdLogger) ConnectFail(m string, err error) {
	msg := err.Error()
	lm := &LogMessage{
		Message: "Connect Fail: " + m,
		Error:   msg,
	}
	logger.Warn(lm)
}

func (l *StdLogger) WillPanic(m string, err error) {
	msg := err.Error()
	lm := &LogMessage{
		Message: m,
		Error:   msg,
	}

	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not marshal will panic LogMessage", lm, err)
	} else {
		logger.Critical(string(bytes))
	}
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
		Message: "Panic Recovered: " + m,
		Error:   errMsg,
		Value:   str,
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not marshal panic LogMessage", lm, err)
	} else {
		logger.Critical(string(bytes))
	}
}

func (l *StdLogger) Error(m string, e error) {
	msg := e.Error()
	lm := &LogMessage{
		Message: m,
		Error:   msg,
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log error because of marshal fail, from error "+m, lm, err)
		return
	} else {
		logger.Error(string(bytes))
	}
}

func (l *StdLogger) Warn(m string, k string, v interface{}) {
	lm := &LogMessage{
		Message: m,
		Key:     k,
		Value:   v,
	}

	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log warn because of marshal fail", lm, err)
		return
	} else {
		logger.Error(string(bytes))
	}
}

// Infom captures a simple message. If you are logging key value pairs,
// use Info(m interface{})
func (l *StdLogger) Inform(m string) {
	logger.Info(&LogMessage{Message: m})
}

// Info logs key value pairs, typically to JSON. Typically using an anonymous struct:
//
//		log.Info(struct{MyKey string}{MyKey:"value to capture"})
func (l *StdLogger) Event(m string, k string, v interface{}) {
	lm := &LogMessage{
		Message: m,
		Key:     k,
		Value:   v,
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log event because of marshal fail", lm, err)
		return
	} else {
		logger.Info(string(bytes))
	}
}

func (l *StdLogger) Debugf(m string, params ...interface{}) {
	lm := &LogMessage{
		Message: fmt.Sprintf(m, params...),
	}
	logger.Debug(lm)
}
