package multi_test

import (
	"strconv"
	"sync/atomic"

	. "github.com/gotgo/fw/multi"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type CountAction struct {
	Count *int32
}

func (a *CountAction) Run(in interface{}) (interface{}, error) {
	atomic.AddInt32(a.Count, 1)
	return in, nil
}

func (a *CountAction) Name() string { return "countAction" }

var _ = Describe("TaskRun", func() {

	var (
		list      []string
		iteration int
		action    *CountAction
		ctx       *DataContext
	)

	BeforeEach(func() {
		iteration++
		list = make([]string, 27)
		for i := 0; i < len(list); i++ {
			list[i], _ = strconv.Unquote(strconv.QuoteRune('a' + rune(1)))
		}
		ctx = &DataContext{}
		action = &CountAction{
			Count: new(int32),
		}
	})

	AfterEach(func() {
	})

	It("should download with single concurrency", func() {
		task := &TaskRun{
			Action:       action,
			Concurrency:  1,
			MaxQueuedIn:  1,
			MaxQueuedOut: 1,
		}

		task.Startup()
		for _, v := range list {
			task.Add(v, ctx)
			<-task.Completed()
		}
		task.Shutdown()

		Expect(int(*action.Count)).To(Equal(len(list)))
	})

	It("should download with multiple concurrency", func() {
		task := &TaskRun{
			Action:       action,
			Concurrency:  4,
			MaxQueuedIn:  8,
			MaxQueuedOut: len(list),
		}

		task.Startup()
		for _, v := range list {
			task.Add(v, ctx)
		}

		task.Shutdown()

		i := 0
		for _ = range task.Completed() {
			i++
		}

		Expect(int(*action.Count)).To(Equal(len(list)))
	})

})
