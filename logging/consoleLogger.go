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

func (l *ConsoleLogger) MarshalFail(m string, err error) {
	lm := &LogMessage{Message: "failed to Marshal object", Error: err.Error()}
	fmt.Println(lm)
}

func (l *ConsoleLogger) UnmarshalFail(m string, err error) {
	lm := &LogMessage{Message: "failed to Unmarshal object", Error: err.Error()}
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

func (l *ConsoleLogger) Warn(m string, k string, v interface{}) {
	lm := &LogMessage{
		Message: m,
		Key:     k,
		Value:   v,
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.UnmarshalFail("failed to warn", err)
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

func (l *ConsoleLogger) Event(m string, k string, v interface{}) {
	lm := &LogMessage{
		Message: m,
		Key:     k,
		Value:   v,
	}
	fmt.Println(lm)
}

func (l *ConsoleLogger) Error(m string, err error) {
	lm := LogMessage{
		Message: m,
		Error:   err.Error(),
	}
	if bytes, err := json.Marshal(lm); err != nil {
		l.UnmarshalFail("could not log error", err)
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
		l.UnmarshalFail("panic message failed to marshal", err)
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
		l.UnmarshalFail("panic message failed to marshal", err)
	} else {
		fmt.Println(string(bytes))
	}
}

func (l *ConsoleLogger) Debugf(m string, params ...interface{}) {
	fmt.Println(&LogMessage{Message: fmt.Sprintf(m, params...)})
}
