package me

import (
	"sync"

	"github.com/gotgo/fw/logging"
)

var globalLogger logging.LoggerBasic
var logMutex sync.Mutex

func log(l logging.LoggerBasic) logging.LoggerBasic {
	if l != nil {
		return l
	} else if globalLogger != nil {
		return globalLogger
	} else {
		return &logging.NoOpLogger{}
	}
}

func SetGlobalLogger(l logging.LoggerBasic) {
	logMutex.Lock()
	globalLogger = l
	logMutex.Unlock()
}

func LogRecoveredPanic(l logging.LoggerBasic, m string, p interface{}, kv ...*logging.KV) {
	//TODO: change interface to include KV
	//TODO: rename to RecoveredPanic
	log(l).HadPanic(m, p)
}

func LogWillPanic(l logging.LoggerBasic, m string, err error, kv ...*logging.KV) {
	log(l).WillPanic(m, err, kv...)
}

func LogError(l logging.LoggerBasic, m string, err error, kv ...*logging.KV) {
	log(l).Error(m, err, kv...)
}

func LogWarn(l logging.LoggerBasic, m string, kv ...*logging.KV) {
	log(l).Warn(m, kv...)
}

func LogInform(l logging.LoggerBasic, m string, kv ...*logging.KV) {
	log(l).Inform(m, kv...)
}

func LogDebug(l logging.LoggerBasic, m string, kv ...*logging.KV) {
	log(l).Debug(m, kv...)
}
