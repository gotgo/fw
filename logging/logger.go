package logging

type Logger interface {
	MarshalFail(m string, err error)
	UnmarshalFail(m string, err error)

	Timeout(m string, err error)
	ConnectFail(m string, err error)

	HadPanic(m string, p interface{})
	WillPanic(m string, err error)
	Error(m string, err error)
	Warn(m string, k string, v interface{})
	// Inform captures a simple message. If you are logging key value pairs,
	// Inform("Server is starting...")
	Inform(m string)
	// Event logs key value pairs, typically to JSON. Typically using an anonymous struct:
	//		log.Event("messageReceived", "message",  msg)
	Event(m string, k string, v interface{})

	Debugf(m string, args ...interface{})

	Log(m *LogMessage)
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
