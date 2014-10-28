package logging

import (
	"encoding/json"
	"fmt"
)

type ConsoleLogger struct {
}

func (l *ConsoleLogger) Log(m *LogMessage) {
	fmt.Println(m)
}

func (l *ConsoleLogger) MarshalFail(m string, obj interface{}, err error) {
	lm := &LogMessage{
		Message: "failed to Marshal object",
		Error:   err.Error(),
		Key:     "object",
		Value:   fmt.Sprintf("%#v", obj),
	}
	fmt.Println(lm)
}

func (l *ConsoleLogger) UnmarshalFail(m string, data []byte, err error) {
	var persistedData []byte

	const arbitraryCutoffSize = 5000
	if len(data) < arbitraryCutoffSize {
		persistedData = data
	}

	//how do we capture the data?
	lm := &LogMessage{
		Message: "failed to Unmarshal object",
		Error:   err.Error(),
		Key:     "rawData",
		Value:   string(persistedData),
	}
	fmt.Println(lm)
}

func (l *ConsoleLogger) Timeout(m string, err error) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
	}
	fmt.Println(lm)
}
func (l *ConsoleLogger) ConnectFail(m string, err error) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
	}
	fmt.Println(lm)
}

func (l *ConsoleLogger) Warn(m string, kv ...*KeyValue) {
	lm := &LogMessage{
		Message: m,
	}
	SetKeyValue(lm, kv...)
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("failed to warn", lm, err)
		return
	} else {
		fmt.Println(string(bytes))
	}
}

// Infom captures a simple message. If you are logging key value pairs,
// use Info(m interface{})
func (l *ConsoleLogger) Inform(m string) {
	fmt.Println(&LogMessage{Message: m})
}

func (l *ConsoleLogger) Event(m string, kv ...*KeyValue) {
	lm := &LogMessage{
		Message: m,
	}
	SetKeyValue(lm, kv...)
	fmt.Println(lm)
}

func (l *ConsoleLogger) Error(m string, err error) {
	lm := LogMessage{
		Message: m,
		Error:   err.Error(),
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("could not log error", lm, err)
		return
	} else {
		fmt.Println(string(bytes))
	}
}

func (l *ConsoleLogger) HadPanic(m string, r interface{}) {
	err, _ := r.(error)
	str, _ := r.(string)

	lm := &LogMessage{
		Message: "Panic Recovered: " + m,
		Error:   err.Error(),
		Value:   str,
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("panic message failed to marshal", lm, err)
		panic(err)
	} else {
		fmt.Println(string(bytes))
		panic(err)
	}
}

func (l *ConsoleLogger) WillPanic(m string, err error) {
	lm := &LogMessage{
		Message: m,
		Error:   err.Error(),
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.MarshalFail("panic message failed to marshal", lm, err)
	} else {
		fmt.Println(string(bytes))
	}
}

func (l *ConsoleLogger) Debugf(m string, params ...interface{}) {
	fmt.Println(&LogMessage{Message: fmt.Sprintf(m, params...)})
}
