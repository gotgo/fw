package multi_test

import (
	"errors"

	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestAction struct {
	ActionCount int
	Err         error
}

func (ta *TestAction) Action(tuple DataFlow, done func()) {
	ta.ActionCount++
	done()
}
func (ta *TestAction) Error() error {
	return ta.Err
}

var _ = Describe("Coordinator", func() {
	It("should succeed", func() {
		concurrency := 2
		retries := 1
		max := 1
		dl := NewCoordinator("pass", concurrency, retries, max)
		dl.Run()

		flow := NewFlow()
		successAction := &TestAction{}
		flow.NewStep("pass", successAction, nil)

		dl.Act([]*Flow{flow})
		//dl.NoMore()

		<-dl.Finished()
		Expect(successAction.ActionCount).To(Equal(1))
		s := <-dl.Success
		Expect(s).ToNot(BeNil())

	})

	It("should error", func() {
		concurrency := 2
		retries := 0
		max := 1
		dl := NewCoordinator("fail", concurrency, retries, max)
		dl.Run()

		flow := NewFlow()

		errorAction := &TestAction{Err: errors.New("test")}
		flow.NewStep("fail", errorAction, nil)

		dl.Act([]*Flow{flow})
		//dl.NoMore()

		<-dl.Finished()

		// 1 execution and 1 retry
		Expect(errorAction.ActionCount).To(Equal(1))
		e := <-dl.Fail
		Expect(e).ToNot(BeNil())
	})

})
