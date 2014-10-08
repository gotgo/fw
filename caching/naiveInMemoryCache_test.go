package caching_test

import (
	. "github.com/krave-n/go/caching"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type cachedData struct {
	Name string
}

var _ = Describe("NaiveInMemoryCache", func() {
	Context("ObjectCache interface", func() {
		cache := NewNaiveInMemoryCache()

		It("Should set", func() {
			cd := &cachedData{"Fred"}
			err := cache.Set("ns", "key", cd)
			Expect(err).Should(BeNil())
		})

		It("Should get", func() {
			cd := new(cachedData)
			miss, err := cache.Get("ns", "key", &cd)
			Expect(miss).ToNot(BeTrue())
			Expect(err).To(BeNil())
		})

	})

})
