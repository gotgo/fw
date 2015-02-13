package me

import (
	"sync"

	. "github.com/gotgo/fw/logging"
)

var globalLogger LoggerBasic
var logMutex sync.Mutex

func log(l LoggerBasic) LoggerBasic {
	if l != nil {
		return l
	} else if globalLogger != nil {
		return globalLogger
	} else {
		return &NoOpLogger{}
	}
}

func SetGlobalLogger(l LoggerBasic) {
	logMutex.Lock()
	globalLogger = l
	logMutex.Unlock()
}

func LogRecoveredPanic(l LoggerBasic, m string, p interface{}, kv ...*KV) {
	//TODO: change interface to include KV
	//TODO: rename to RecoveredPanic
	log(l).HadPanic(m, p)
}

func LogWillPanic(l LoggerBasic, m string, err error, kv ...*KV) {
	log(l).WillPanic(m, err, kv...)
}

func LogError(l LoggerBasic, m string, err error, kv ...*KV) {
	log(l).Error(m, err, kv...)
}

func LogWarn(l LoggerBasic, m string, kv ...*KV) {
	log(l).Warn(m, kv...)
}

func LogInform(l LoggerBasic, m string, kv ...*KV) {
	log(l).Inform(m, kv...)
}

func LogDebug(l LoggerBasic, m string, kv ...*KV) {
	log(l).Debug(m, kv...)
}
