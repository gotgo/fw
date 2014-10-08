package tracing_test

import (
	"time"

	"github.com/gotgo/fw/logging"
	. "github.com/krave-n/go/tracing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockReceiver struct {
	ReceivedCount int
	Closed        bool
	Delay         time.Duration
}

func (mr *mockReceiver) Name() string {
	return "mock"
}

func (mr *mockReceiver) Receive(tm *TraceMessage) {
	mr.ReceivedCount++
	if mr.Delay > 0 {
		time.Sleep(mr.Delay)
	}
}

func (mr *mockReceiver) Close() {
	mr.Closed = true
}

var _ = Describe("Dispatcher", func() {

	var (
		dispatcher *Dispatcher
		aMessage   *TraceMessage
	)

	BeforeEach(func() {
		dispatcher = &Dispatcher{
			Log: new(logging.NoOpLogger),
		}

		aMessage = &TraceMessage{}
		dispatcher.Start()
	})

	AfterEach(func() {
		dispatcher.Stop()
	})

	It("should gracefully handle multiple calls to Start", func() {
		dispatcher.Start()
		dispatcher.Start()
	})

	It("should gracefully handle multiple calls to Stop", func() {
		dispatcher.Stop()
		dispatcher.Stop()
	})

	It("should gracefully handle calls to Capture after Stop", func() {
		dispatcher.Stop()
		dispatcher.Capture(aMessage)
	})

	It("should close all receivers, on close", func() {
		r1 := &mockReceiver{}
		r2 := &mockReceiver{}
		dispatcher.Register(r1)
		dispatcher.Register(r2)
		dispatcher.Stop()
		Expect(r1.Closed).To(BeTrue())
		Expect(r2.Closed).To(BeTrue())
	})

	It("should increment dropped messages if Capture is called after drop", func() {
		dispatcher.Stop()
		dispatcher.Capture(aMessage)
		Expect(dispatcher.DroppedMessages()).To(Equal(int64(1)))
	})

	Context("multi-threading", func() {
		It("should send a Captured Message to all receivers", func() {
			r1 := &mockReceiver{}
			r2 := &mockReceiver{}
			dispatcher.Register(r1)
			dispatcher.Capture(aMessage)
			//async
			time.Sleep(100 * time.Millisecond)
			Expect(r1.ReceivedCount).To(Equal(1))

			dispatcher.Register(r2)
			dispatcher.Capture(aMessage)
			//async
			time.Sleep(100 * time.Millisecond)
			Expect(r1.ReceivedCount).To(Equal(2))
			Expect(r2.ReceivedCount).To(Equal(1))
		})

		It("should not wait for messages to finish on stop", func() {
			m := &mockReceiver{
				Delay: 100 * time.Millisecond,
			}
			count := 6
			dispatcher.Register(m)

			for i := 0; i < count; i++ {
				dispatcher.Capture(aMessage)
			}

			Expect(dispatcher.PendingMessages()).To(Equal(count))
			dispatcher.Stop()
			Expect(m.ReceivedCount).To(BeNumerically("<", count))
		})

		It("should drop messages if there are too many queued", func() {
			m := &mockReceiver{
				Delay: 10 * time.Millisecond,
			}
			dispatcher.Register(m)
			count := dispatcher.QueueSize() + 200

			for i := 0; i < count; i++ {
				dispatcher.Capture(aMessage)
			}

			dispatcher.Stop()
			Expect(dispatcher.DroppedMessages()).To(BeNumerically(">", 0))
		})

		It("should not close a writing receiver", func() {
			m := &mockReceiver{
				Delay: 500 * time.Millisecond,
			}

			dispatcher.Register(m)
			dispatcher.Capture(aMessage)
			time.Sleep(50 * time.Millisecond)
			start := time.Now()
			dispatcher.Stop()
			took := time.Since(start)
			Expect(m.Closed).To(BeTrue())
			Expect(took).To(BeNumerically(">", 200*time.Millisecond))
		})
	})
})
