package logging

type Logger interface {
	LoggerBasic
	// MarshalFail occurs when an object fails to marshal.
	// Solving a Marshal failure requires discovering which object type and what data was
	// in that instance that could have caused the failure. This is why the interface requires
	// the object
	MarshalFail(m string, obj interface{}, err error)
	// UnmarshalFail occures when a stream is unable to be unmarshalled.
	// Solving a unmarshal failure requires knowing what object type, which field, and
	// what's wrong with the source data that causes the problem
	UnmarshalFail(m string, data []byte, err error)

	Timeout(m string, err error, kv ...*KV)
	ConnectFail(m string, err error, kv ...*KV)
}

type LoggerBasic interface {
	HadPanic(m string, p interface{})
	WillPanic(m string, err error, kv ...*KV)

	Error(m string, err error, kv ...*KV)
	Warn(m string, kv ...*KV)

	// Inform captures a simple message.
	// Inform("Server is starting...")
	Inform(m string, kv ...*KV)

	Event(m string, kv ...*KV)
	Debug(m string, kv ...*KV)

	Log(m *LogMessage)
}
