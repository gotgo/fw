package tracing

import "io"

//type RequestTracer interface {
//	Encode() func(v interface{}) ([]byte, error)
//}

type Content struct {
	Request  interface{} `json:"request,omitempty"`
	Response interface{} `json:"response,omitempty"`
	Args     interface{} `json:"args,omitempty"`
	Headers  interface{} `json:"headers,omitempty"`
}

type CallerIdentityTrace struct {
	AccountUid  string `json:"acctUid,omitempty"`
	LocationUid string `json:"locationUid, omitempty"`
	DeviceUid   string `json:"deviceUid, omitempty"`
}

type TraceMessageWriter interface {
	Write(b []byte) (n int, err error)
	Close() error
}

type ContentType string

const (
	Json ContentType = "application/json"
)

type From string

const (
	Panic          From = "panic"
	Error          From = "error"
	RequestTimeout From = "reqTimeout"
	ConnectTimeout From = "conTimeout"
	IOError        From = "ioError"
	RequestData    From = "requestData"
	ResponseData   From = "responseData"
	MarshalError   From = "marshalError"
	UnmarshalError From = "unmarshalError"
)

type Tracer interface {
	Annotate(f From, k string, v interface{})
	AnnotateBinary(f From, k string, reader io.Reader, ct string)
	NewRequest(name string, args interface{}, headers map[string][]string) RequestTracer
}

type RequestTracer interface {
	Annotate(f From, k string, v interface{})
	AnnotateBinary(f From, k string, reader io.Reader, ct string)
	Begin()
	End()
}
