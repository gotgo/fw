package tracing

type TraceReceiver interface {
	Name() string
	Receive(*TraceMessage)
	Close()
}
