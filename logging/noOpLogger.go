package logging

type NoOpLogger struct {
}

func (l *NoOpLogger) Log(m *LogMessage)                 {}
func (l *NoOpLogger) MarshalFail(m string, err error)   {}
func (l *NoOpLogger) UnmarshalFail(m string, err error) {}

func (l *NoOpLogger) Timeout(m string, err error)     {}
func (l *NoOpLogger) ConnectFail(m string, err error) {}

func (l *NoOpLogger) WillPanic(m string, err error)          {}
func (l *NoOpLogger) HadPanic(m string, r interface{})       {}
func (l *NoOpLogger) Error(m string, err error)              {}
func (l *NoOpLogger) Warn(m string, k string, v interface{}) {}

func (l *NoOpLogger) Inform(m string)                         {}
func (l *NoOpLogger) Event(m string, k string, v interface{}) {}

func (l *NoOpLogger) Debugf(m string, params ...interface{}) {}
