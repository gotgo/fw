package tracing

import (
	"sync"
	"sync/atomic"

	"github.com/gotgo/fw/logging"
)

const channelSize = 2000

type Dispatcher struct {
	Log             logging.Logger
	disabled        bool
	messages        chan *TraceMessage
	receivers       []TraceReceiver
	isRunning       bool
	mutex           sync.Mutex
	stop            chan struct{}
	remaining       *sync.WaitGroup
	messagesPending int32
	droppedMessages int64
}

func (d *Dispatcher) QueueSize() int {
	return channelSize
}

func (d *Dispatcher) DroppedMessages() int64 {
	return d.droppedMessages
}

func (d *Dispatcher) PendingMessages() int {
	return int(d.messagesPending)
}

func (d *Dispatcher) Start() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.isRunning == false {
		d.log("Starting trace dispatcher.")

		d.messages = make(chan *TraceMessage, channelSize)
		d.stop = make(chan struct{})
		d.remaining = new(sync.WaitGroup)
		d.remaining.Add(1)

		d.messagesPending = 0
		d.droppedMessages = 0
		go d.runAsync()
		d.isRunning = true
	}
}

func (d *Dispatcher) Register(receiver TraceReceiver) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.receivers = append(d.receivers, receiver)
	d.log("Registered trace receiver: " + receiver.Name())
}

func (d *Dispatcher) Stop() {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if d.isRunning {
		d.log("Stopping Trace dispatcher...")
		close(d.stop)
		close(d.messages)

		d.remaining.Wait()
		d.isRunning = false
		d.messagesPending = 0
		d.log("Trace Dispatcher Stopped")
	}
}

func (d *Dispatcher) Capture(message *TraceMessage) {
	//nil channel blocks
	//closed channel receives
	if d.isRunning {
		if count := atomic.AddInt32(&d.messagesPending, 1); count >= channelSize {
			dropped := atomic.AddInt64(&d.droppedMessages, 1)
			d.Log.Warn("dropped trace message.", &logging.KeyValue{"droppedCount", dropped})
			return
		}
		d.messages <- message
	} else {
		atomic.AddInt64(&d.droppedMessages, 1)
	}
}

func (d *Dispatcher) log(message string) {
	if log := d.Log; log != nil {
		log.Inform(message)
	}
}

func (d *Dispatcher) runAsync() {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			msg, _ := r.(string)
			d.Log.Error("panic in dispatcher "+msg, err)
		}
		d.remaining.Done()
	}()

	for {
		select {
		case m := <-d.messages:
			rec := d.receivers
			for i := 0; i < len(rec); i++ {
				rec[i].Receive(m)
			}
			atomic.AddInt32(&d.messagesPending, -1)
		case <-d.stop:
			rec := d.receivers
			for i := 0; i < len(rec); i++ {
				r := rec[i]
				d.log("Closing trace receiver: " + r.Name())
				r.Close()
			}
			return
		}
	}
}
