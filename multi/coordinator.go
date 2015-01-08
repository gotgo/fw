package multi

import (
	"sync"
	"sync/atomic"

	"github.com/gotgo/fw/me"
)

//  NewCoordintor
func NewCoordinator(name string, concurrency int, retries int, maxItems int) *Coordinator {
	size := maxItems
	d := &Coordinator{
		todo: make(chan *Flow, size),
		//	retry: make(chan *Flow, size),
		//rateLimiter: make(chan interface{}, concurrency),
		//		completed: make(chan *Flow, size),
		//stop:     make(chan interface{}, size),
		finished: make(chan struct{}, size),
		Success:  make(chan *Flow, size),
		Fail:     make(chan *Flow, size),
		queued:   new(int32),
		name:     name,
		retries:  retries,
		maxSize:  size,
	}
	return d
}

// Coordinator - concurrent actions
type Coordinator struct {
	todo chan *Flow
	//rateLimiter chan interface{}
	//	completed chan *Flow
	//retry     chan *Flow

	finished chan struct{}
	//stop     chan interface{}

	Success chan *Flow
	Fail    chan *Flow

	name    string
	retries int

	queued  *int32
	maxSize int
}

func (d *Coordinator) Finished() <-chan struct{} {
	return d.finished
}

func (d *Coordinator) Run() {
	go d.feedTodo()
	//go d.feedRetry()
	//	go d.handleResults()
}

func (c *Coordinator) noMore() {
	close(c.todo)
	//	c.stop <- nil
}

func (c *Coordinator) isComplete() bool {
	return *c.queued == int32(0)
}

func (d *Coordinator) Act(flows []*Flow) {
	//	if len(flows) > d.maxSize {
	//		panic("Act() the lengths of the flows can not be greater than the max")
	//	}
	for _, f := range flows {
		atomic.AddInt32(d.queued, 1)
		d.todo <- f
	}
	d.noMore()
}

func (d *Coordinator) From(c *Coordinator) {
	go d.from(c)
}

func (d *Coordinator) from(c *Coordinator) {
	for f := range c.Success {
		atomic.AddInt32(d.queued, 1)
		d.todo <- f
	}
	d.noMore()
}

func (c *Coordinator) process(f *Flow, done sync.WaitGroup) {
	currentStep := f.Steps[c.name]
	if currentStep == nil {
		panic(me.NewErr("unknown step " + c.name))
	}
	currentStep.Action.Action(f.Data, func() {
		if currentStep.Action.Error() != nil {
			c.Fail <- f
		} else {
			c.Success <- f
		}
		atomic.AddInt32(c.queued, -1)
		//		<-c.rateLimiter //remove one, any one, to allow more
		done.Done()
	})
	//	c.rateLimiter <- nil //rate limited, will block at limit
}

func (c *Coordinator) feedTodo() {
	var done sync.WaitGroup
	for f := range c.todo {
		done.Add(1)
		go c.process(f, done) //synchronous
	}
	done.Wait()
	c.closeUp()
}

//func (c *Coordinator) feedRetry() {
//	for f := range c.retry {
//		c.process(f)
//	}
//}

//func (d *Coordinator) handleResults() {
//	shutdown := false
//	for {
//		select {
//		case f := <-d.completed:
//			<-d.rateLimiter //remove one, any one, to allow more
//			s := f.Steps[d.name]
//			s.Attempts++
//			flow := f //local
//
//			if s.Action.Error() != nil {
//				//if s.Attempts > d.retries {
//				d.Fail <- flow
//				atomic.AddInt32(d.queued, -1)
//				//} else {
//				//	d.retry <- flow
//				//	fmt.Println("retry")
//				//}
//			} else {
//				d.Success <- flow
//				atomic.AddInt32(d.queued, -1)
//			}
//
//			if shutdown && d.isComplete() {
//				close(d.stop)
//			}
//		case <-d.stop:
//			if d.isComplete() {
//				d.closeUp()
//				return
//			} else {
//				shutdown = true
//			}
//		}
//	}
//}

func (c *Coordinator) closeUp() {
	close(c.finished)
	close(c.Success)
	close(c.Fail)
	//	close(c.retry)
}
