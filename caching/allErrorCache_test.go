package caching_test

import (
	. "github.com/krave-n/go/caching"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AllErrorCache", func() {
	Context("ObjectCache", func() {
		It("should return errors", func() {
			ec := new(AllErrorCache)
			_, err := ec.Get("", "", nil)
			Expect(err).ToNot(BeNil())
			err = ec.Set("", "", nil)
			Expect(err).ToNot(BeNil())
		})
	})

})
