package zoo_test

import (
	. "github.com/gotgo/fw/zoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KafkaState", func() {
	It("should create the correct path", func() {
		k := &KafkaState{
			Root:  "a",
			Topic: "mytopic",
		}
		Expect(k.PartitionsPath()).To(Equal("/a/brokers/topics/mytopic/partitions"))
	})
})
