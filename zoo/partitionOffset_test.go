package zoo_test

import (
	"sort"

	. "github.com/gotgo/fw/zoo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PartitionOffset", func() {

	It("should work", func() {
		p := []*PartitionOffset{
			{5, 1},
			{2, 1},
			{4, 1},
			{1, 1},
		}

		sort.Sort(ByPartition{p})
		Expect(p[0].Partition).To(Equal(int32(1)))
		Expect(p[1].Partition).To(Equal(int32(2)))
		Expect(p[2].Partition).To(Equal(int32(4)))
		Expect(p[3].Partition).To(Equal(int32(5)))
	})

})
