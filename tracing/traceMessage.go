package tracing

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"time"

	"github.com/satori/go.uuid"
)

// RequestTrace
//
//	To Consider:
//  Events == Annotations
//  Identity (UserId, DeviceId, AccountId, LocationId)
type TraceMessage struct {
	Duration    time.Duration `json:"duration"`
	BeganAt     time.Time     `json:"beganAt"`
	Name        string        `json:"name"`
	TraceUid    string        `json:"traceUid"`
	SpanUid     string        `json:"spanUid"`
	PrevSpanUid string        `json:"prevSpanUid,omitempty"`
	Role        string        `json:"role"`              //sender, receiver
	Outcome     string        `json:"outcome,omitempty"` //ok, error
	Content     *Content      `json:"content,omitempty"`
	Host        *Host         `json:"host"`
	Annotations []*Annotation `json:"annotations"`
}

func NewReceiveTrace(traceUid, existingSpanUid string) *TraceMessage {
	if traceUid == "" {
		traceUid = newUid()
	}
	if existingSpanUid == "" {
		existingSpanUid = newUid()
	}

	trace := &TraceMessage{
		TraceUid: traceUid,
		SpanUid:  existingSpanUid,
		Role:     "receive",
	}
	return trace
}

func NewRequestTrace(traceUid, prevSpanUid string) *TraceMessage {
	if traceUid == "" {
		traceUid = newUid()
	}

	trace := &TraceMessage{
		TraceUid:    traceUid,
		SpanUid:     newUid(),
		PrevSpanUid: prevSpanUid,
		Role:        "request",
	}
	return trace
}

func (tm *TraceMessage) NewRequest(name string, args interface{}, headers map[string][]string) RequestTracer {
	msg := NewRequestTrace(tm.TraceUid, tm.SpanUid)
	msg.Name = name
	msg.Content = &Content{Args: args, Headers: headers}
	return msg
}

func (tm *TraceMessage) Begin() {
	tm.BeganAt = time.Now()
}

func (tm *TraceMessage) End() {
	tm.RequestCompleted()
}

func (tm *TraceMessage) Annotate(f From, k string, v interface{}) {
	tm.Annotations = append(tm.Annotations, &Annotation{From: f, Name: k, Value: v})
}

func (tm *TraceMessage) AnnotateBinary(f From, k string, reader io.Reader, ct string) {
	var data interface{}
	bts, err := ioutil.ReadAll(reader)
	if err != nil {
		tm.Annotate(f, k, "binary annotate failed "+err.Error())
	}

	if ct == "application/json" {
		if err := json.Unmarshal(bts, &data); err != nil {
			m := data.(map[string]interface{})
			tm.Annotations = append(tm.Annotations, &Annotation{From: f, Name: k, Value: m, ContentType: ct, IsBinary: false})
		}
	} else {
		tm.Annotations = append(tm.Annotations, &Annotation{From: f, Name: k, Value: bts, ContentType: ct, IsBinary: true})
	}
}

func (tm *TraceMessage) Error(name string, value interface{}) {
	tm.Annotations = append(tm.Annotations, &Annotation{Name: name, Value: value, From: Error})
}

func (tm *TraceMessage) RequestCompleted() {
	elapsed := time.Since(tm.BeganAt)
	tm.Duration = elapsed
}

func (tm *TraceMessage) RequestFail() {
	elapsed := time.Since(tm.BeganAt)
	tm.Outcome = "error"
	tm.Duration = elapsed
}

func newUid() string {
	bytes := uuid.NewV4().Bytes()
	return base64.StdEncoding.EncodeToString(bytes)
}

func (tm *TraceMessage) ReceivedRequest(requestName string, args map[string]string, headers map[string][]string) {
	tm.Name = requestName
	tm.BeganAt = time.Now()
	tm.Role = "received"
	//TODO: filter sensitive header info
	tm.Content = &Content{Args: args, Headers: headers}
}
