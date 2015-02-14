package multi_test

import (
	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RecentlyExisted", func() {

	It("should work", func() {
		e := NewRecentlyExisted(3)
		r := e.CheckAndAdd("a")
		Expect(r).To(BeFalse())
		r = e.CheckAndAdd("b")
		Expect(r).To(BeFalse())
		r = e.CheckAndAdd("a")
		Expect(r).To(BeTrue())
		r = e.CheckAndAdd("c")
		Expect(r).To(BeFalse())
		r = e.CheckAndAdd("d")
		Expect(r).To(BeFalse())
		r = e.CheckAndAdd("a")
		Expect(r).To(BeFalse())
	})

})
