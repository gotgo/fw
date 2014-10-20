package caching_test

import (
	. "github.com/gotgo/fw/caching"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NoOpCache", func() {

	It("should do nothing", func() {
		cache := new(NoOpCache)
		err := cache.Set("ns", "key", struct{ Name string }{"Fred"})
		Expect(err).To(BeNil())

		var data = new(struct{ Name string })
		miss, err := cache.Get("ns", "key", &data)
		Expect(err).To(BeNil())
		Expect(miss).To(BeTrue())
	})
})
