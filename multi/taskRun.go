package multi

import (
	"sync"

	"github.com/gotgo/fw/logging"
	"github.com/gotgo/fw/me"
	"github.com/gotgo/fw/stats"
)

func (t *TaskRun) setup() {
	if t.Concurrency == 0 {
		t.Concurrency = 1
	}
	if t.MaxQueuedIn == 0 {
		t.MaxQueuedIn = 1
	}
	if t.MaxQueuedOut == 0 {
		t.MaxQueuedOut = 1
	}
	t.input = make(chan *TaskRunInput, t.MaxQueuedIn)
	t.output = make(chan *TaskRunOutput, t.MaxQueuedOut)
	t.shutdown = make(chan struct{})
	t.done = make(chan struct{})
	t.outstanding = &sync.WaitGroup{}
	t.closeDoneOnce = &sync.Once{}
	t.closeInputOnce = &sync.Once{}
	//t.Track = stats.NewBasicMeter("image.downloader", me.App.Environment())
}

func (t *TaskRun) Startup() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if t.running {
		panic("already running")
	}
	t.running = true
	t.setup()
	concurrency := t.Concurrency

	for i := 0; i < concurrency; i++ {
		go t.run()
	}
}

type TaskRun struct {
	Action TaskAction

	Config map[string]interface{}

	mutex          sync.Mutex
	Track          stats.BasicMeter
	Log            logging.Logger
	Concurrency    int
	MaxQueuedIn    int
	MaxQueuedOut   int
	running        bool
	input          chan *TaskRunInput
	output         chan *TaskRunOutput
	outstanding    *sync.WaitGroup
	closeDoneOnce  *sync.Once
	closeInputOnce *sync.Once
	done           chan struct{}
	shutdown       chan struct{}
}

type TaskRunResult struct {
	Error  error
	Input  interface{}
	Output interface{}
}

type TaskRunInput struct {
	Input   interface{}
	Context *DataContext
}

type TaskRunOutput struct {
	result  *TaskRunResult
	Context *DataContext
}

func (o *TaskRunOutput) Error() error {
	return o.result.Error
}
func (o *TaskRunOutput) Input() interface{} {
	return o.result.Input
}
func (o *TaskRunOutput) Output() interface{} {
	return o.result.Output
}

func (o *TaskRunOutput) Previous(name string) *TaskRunResult {
	return o.Context.Get(name).(*TaskRunResult)
}

// Add - will block when the number of items queued reaches MaxQueuedInput
func (t *TaskRun) Add(todo interface{}, context *DataContext) {
	t.input <- &TaskRunInput{Input: todo, Context: context}
}

func (t *TaskRun) Completed() <-chan *TaskRunOutput {
	return t.output
}

// Shutdown - begins the shutdown operation. Reading on the returned channel will block until the shutdown is complete
func (t *TaskRun) Shutdown() chan struct{} {
	//no more input
	t.closeInputOnce.Do(func() { close(t.input) })
	return t.done
}

func (t *TaskRun) Name() string {
	return t.Action.Name()
}

func (t *TaskRun) run() {
	for in := range t.input {
		t.safeExecute(in)
	}
	t.outstanding.Wait()

	t.closeDoneOnce.Do(func() {
		close(t.output)
		close(t.done)
	})
}

func (t *TaskRun) safeExecute(task *TaskRunInput) {
	t.outstanding.Add(1)
	defer func() {
		t.outstanding.Done()
		if r := recover(); r != nil {
			r := &TaskRunResult{
				Error: me.NewErr("recovered from panic"),
			}
			task.Context.Set(t.Name(), r)
			t.output <- &TaskRunOutput{
				result:  r,
				Context: task.Context,
			}
			me.LogRecoveredPanic(t.Log, "execute failed", r, &logging.KV{"from", t})
		}
	}()

	out, err := t.Action.Run(task.Input)

	result := &TaskRunResult{
		Error:  err,
		Input:  task.Input,
		Output: out,
	}

	task.Context.Set(t.Name(), result)

	t.output <- &TaskRunOutput{
		result:  result,
		Context: task.Context,
	}
}
