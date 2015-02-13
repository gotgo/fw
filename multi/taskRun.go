package multi

import (
	"sync"

	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
	"github.com/gotgo/fw/stats"
)

func (d *TaskRun) setup() {
	if d.Concurrency == 0 {
		d.Concurrency = 1
	}
	if d.MaxQueuedIn == 0 {
		d.MaxQueuedIn = 1
	}
	if d.MaxQueuedOut == 0 {
		d.MaxQueuedOut = 1
	}
	d.input = make(chan *TaskRunInput, d.MaxQueuedIn)
	d.output = make(chan *TaskRunResult, d.MaxQueuedOut)
	d.shutdown = make(chan struct{})
	d.outstanding = &sync.WaitGroup{}
	d.once = &sync.Once{}
	//d.Track = stats.NewBasicMeter("image.downloader", me.App.Environment())
}

func (d *TaskRun) Start() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.running {
		panic("already running")
	}
	d.running = true
	d.setup()
	concurrency := d.Concurrency

	for i := 0; i < concurrency; i++ {
		go d.run()
	}
}

type TaskRun struct {
	Action TaskAction // func(interface{}) (interface{}, error)

	Config map[string]interface{}

	mutex        sync.Mutex
	Track        stats.BasicMeter
	Log          logging.Logger
	Concurrency  int
	MaxQueuedIn  int
	MaxQueuedOut int
	running      bool
	input        chan *TaskRunInput
	output       chan *TaskRunResult
	outstanding  *sync.WaitGroup
	once         *sync.Once
	done         chan struct{}
	shutdown     chan struct{}
}

type TaskRunResult struct {
	Error   error
	Input   interface{}
	Output  interface{}
	Context *DataContext
}

type TaskRunInput struct {
	Input   interface{}
	Context *DataContext
}

// Add - will block when the number of items queued reaches MaxQueuedInput
func (d *TaskRun) Add(todo interface{}, context *DataContext) {
	d.input <- &TaskRunInput{Input: todo, Context: context}
}

func (d *TaskRun) Completed() <-chan *TaskRunResult {
	return d.output
}

// Shutdown - begins the shutdown operation. Reading on the returned channel will block until the shutdown is complete
func (d *TaskRun) Shutdown() chan struct{} {
	close(d.input) //no more input
	return d.done
}

func (d *TaskRun) run() {
	for in := range d.input {
		d.safeExecute(in)
	}
	d.outstanding.Wait()

	d.once.Do(func() {
		close(d.output)
		close(d.done)
	})
}

func (d *TaskRun) safeExecute(task *TaskRunInput) {
	d.outstanding.Add(1)
	defer func() {
		d.outstanding.Done()
		if r := recover(); r != nil {
			me.LogRecoveredPanic(d.Log, "download failed", r, &logging.KV{"from", d})
		}
	}()

	out, err := d.Action.Run(task.Input)
	d.output <- &TaskRunResult{
		Error:   err,
		Input:   task.Input,
		Output:  out,
		Context: task.Context,
	}
}
