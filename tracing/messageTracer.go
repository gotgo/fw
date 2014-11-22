package tracing

import "io"

type MessageTracer struct {
	message *TraceMessage
}

func NewMessageTracer(message *TraceMessage) *MessageTracer {
	return &MessageTracer{
		message: message,
	}
}

func (mt *MessageTracer) NewRequest(name string, args interface{}, headers map[string][]string) RequestTracer {
	msg := NewRequestTrace(mt.message.TraceUid, mt.message.SpanUid)
	msg.Name = name
	msg.Content = &Content{Args: args, Headers: headers}
	return msg
}

func (td *MessageTracer) Annotate(f From, k string, v interface{}) {
	td.message.Annotate(f, k, v)
}

func (td *MessageTracer) AnnotateBinary(f From, k string, reader io.Reader, ct string) {
	td.message.AnnotateBinary(f, k, reader, ct)
}

func (td *MessageTracer) Error(name string, value error) {
	td.message.Error(name, value)
}

func (td *MessageTracer) TraceUid() string {
	return td.message.TraceUid
}

func (td *MessageTracer) SpanUid() string {
	return td.message.SpanUid
}
