package tracing_test

import (
	"bytes"
	"encoding/json"

	"github.com/gotgo/fw/tracing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TraceData struct {
	Message string `json:"message"`
}

var _ = Describe("TraceMessage", func() {

	var msg *tracing.TraceMessage

	BeforeEach(func() {
		msg = &tracing.TraceMessage{
			Name: "test",
		}
	})

	PContext("Binary annotations", func() {
		It("should capture json bytes", func() {
			td := &TraceData{
				Message: "dude",
			}
			bts, _ := json.Marshal(td)
			msg.AnnotateBinary("test", "test", bytes.NewReader(bts), "application/json")
			//this is throwing
			a := msg.Annotations[0].Value.(map[string]interface{})
			Expect(td.Message).To(Equal(a["message"]))
		})

	})
})
