package kafka_test

import (
	. "github.com/gotgo/fw/kafka"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KafkaClient", func() {

	Context("Split By", func() {
		It("should split partitions equally if partitions is a multiple of splitBy", func() {
			partitions := 5
			splitBy := 3

			buckets := DividePartitions(partitions, splitBy)
			Expect(len(buckets)).To(Equal(splitBy))
			Expect(len(buckets[0])).To(Equal(2))
			Expect(len(buckets[1])).To(Equal(2))
			Expect(len(buckets[2])).To(Equal(1))

		})

		It("should split partitions by the partition size when the partitions is less than the splitBy", func() {
			partitions := 2
			splitBy := 3

			buckets := DividePartitions(partitions, splitBy)
			Expect(len(buckets)).To(Equal(splitBy))
			Expect(len(buckets[0])).To(Equal(1))
			Expect(len(buckets[1])).To(Equal(1))
			Expect(len(buckets[2])).To(Equal(0))
		})
	})
})
