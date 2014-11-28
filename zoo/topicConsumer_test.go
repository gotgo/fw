package zoo_test

import (
	. "github.com/gotgo/fw/zoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TopicConsumer", func() {

	It("should create the correct path", func() {
		t := &TopicConsumer{
			Root:  "a",
			Topic: "mytopic",
			App:   "myapp",
		}

		path := t.PartitionsPath()
		Expect(path).To(Equal("/a/mytopic/apps/myapp/partitions"))
	})
})
