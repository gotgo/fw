package multi

import (
	"sync/atomic"
	"time"
)

//  NewCoordintor
func NewCoordinator(name string, concurrency int, retries int) *Coordinator {
	d := &Coordinator{
		todo:        make(chan *Flow),
		rateLimiter: make(chan interface{}, concurrency),
		completed:   make(chan *Flow),
		stop:        make(chan interface{}),
		finished:    make(chan struct{}),
		in:          new(int32),
		out:         new(int32),
		name:        name,
		retries:     retries,
	}
	return d
}

// Coordinator - concurrent actions
type Coordinator struct {
	todo        chan *Flow
	rateLimiter chan interface{}
	completed   chan *Flow

	finished chan struct{}
	stop     chan interface{}

	Success chan *Flow
	Fail    chan *Flow

	name    string
	retries int

	in  *int32
	out *int32
}

func (d *Coordinator) Finished() <-chan struct{} {
	return d.finished
}

func (d *Coordinator) Run() {
	go d.feed()
	go d.process()
}

func (d *Coordinator) NoMore() {
	close(d.todo)
	d.stop <- nil
}

func (d *Coordinator) isComplete() bool {
	return d.in == d.out
}

func (d *Coordinator) Act(flows []*Flow) {
	for _, f := range flows {
		atomic.AddInt32(d.in, 1)
		d.todo <- f
	}
}

func (d *Coordinator) From(c *Coordinator) {
	go d.from(c)
}

func (d *Coordinator) from(c *Coordinator) {
	for f := range c.Success {
		atomic.AddInt32(d.in, 1)
		d.todo <- f
	}
	d.NoMore()
}

func (d *Coordinator) feed() {
	for f := range d.todo {
		d.rateLimiter <- nil //rate limited, will block at limit
		s := f.Steps[d.name].Source
		go s.Action(func() { d.completed <- f })
	}
}

func (d *Coordinator) process() {
	for {
		select {
		case f := <-d.completed:
			<-d.rateLimiter //remove one, any one, to allow more
			s := f.Steps[d.name]
			s.Attempts++

			if s.Source.Error() != nil {
				if s.Attempts > 3 {
					d.Fail <- f
					atomic.AddInt32(d.out, 1)
				} else {
					d.todo <- f
				}
			} else {
				d.Success <- f
				atomic.AddInt32(d.out, 1)
			}
		case <-d.stop:
			if d.isComplete() {
				close(d.finished)
				close(d.Success)
				close(d.Fail)
				return
			} else {
				go func() { d.stop <- time.After(time.Second) }()
			}
		}
	}
}
