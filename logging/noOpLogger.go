package logging

func NewNoOpLogger() *NoOpLogger {
	return noOpLogger
}

var noOpLogger *NoOpLogger

func init() {
	noOpLogger = new(NoOpLogger)
}

type NoOpLogger struct {
}

func (l *NoOpLogger) Log(m *LogMessage)                                {}
func (l *NoOpLogger) MarshalFail(m string, obj interface{}, err error) {}
func (l *NoOpLogger) UnmarshalFail(m string, data []byte, err error)   {}

func (l *NoOpLogger) Timeout(m string, err error, kv ...*KeyValue)     {}
func (l *NoOpLogger) ConnectFail(m string, err error, kv ...*KeyValue) {}

func (l *NoOpLogger) HadPanic(m string, r interface{})               {}
func (l *NoOpLogger) WillPanic(m string, err error, kv ...*KeyValue) {}
func (l *NoOpLogger) Error(m string, err error, kv ...*KeyValue)     {}
func (l *NoOpLogger) Warn(m string, kv ...*KeyValue)                 {}

func (l *NoOpLogger) Inform(m string)                 {}
func (l *NoOpLogger) Event(m string, kv ...*KeyValue) {}

func (l *NoOpLogger) Debug(m string, kv ...*KeyValue) {}
