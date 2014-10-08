package tracing

import "io"

type NopTracer struct{}

func (td *NopTracer) Annotate(f From, k string, v interface{})                      {}
func (not *NopTracer) AnnotateBinary(f From, k string, reader io.Reader, ct string) {}

func (td *NopTracer) NewRequest(name string, args interface{}, headers map[string][]string) RequestTracer {
	return new(NopClientTracer)
}

type NopClientTracer struct{}

func (nct *NopClientTracer) Annotate(f From, k string, v interface{})                     {}
func (nct *NopClientTracer) AnnotateBinary(f From, k string, reader io.Reader, ct string) {}
func (nct *NopClientTracer) Begin()                                                       {}
func (nct *NopClientTracer) End()                                                         {}
