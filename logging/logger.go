package logging

type Logger interface {

	// MarshalFail occurs when an object fails to marshal.
	// Solving a Marshal failure requires discovering which object type and what data was
	// in that instance that could have caused the failure. This is why the interface requires
	// the object
	MarshalFail(m string, obj interface{}, err error)
	// UnmarshalFail occures when a stream is unable to be unmarshalled.
	// Solving a unmarshal failure requires knowing what object type, which field, and
	// what's wrong with the source data that causes the problem
	UnmarshalFail(m string, data []byte, err error)

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
