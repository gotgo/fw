package stats

import (
	"os"
	"time"
)

func NewBasicMeter(name, environment string, tags ...string) BasicMeter {
	hostName, _ := os.Hostname()
	return &Metrics{
		Hostname:    hostName,
		Name:        name,
		Environment: environment,
		Tags:        tags,
	}
}

// BasicMeter
// Discussion Points - Why does duration & size exist on the interface when they can be
// covered by the other interfaces?
// A: Believe that fundamental building blocks, the atoms, of a system should be easily tracked
// Duration and Data Size are fundamental building blocks for reasoning about a software system.
type BasicMeter interface {
	// Occurence - Something worth measuring happened
	// Data worth knowing:
	// How many times total did it happen
	// How many event / time is happening
	Occurence(name string) //moment

	// Duration - Called after a process has completed, captures both that an
	// occurence of a process happened and how long it took.
	Duration(name string, start time.Time)

	// Distribution - Capture the distribution of a value. Mean, Max, Min, Avg, 95%
	Distribution(name string, value float64)

	// Current - A value where the actual value is the only thing that matters. A level or gauge.
	Value(name string, value float64)

	// Size - Track payload sizes
	Size(name string, size int64)
}

type OutcomeMeter interface {
	// Error - this is not meant to log the error, just count that an occur occured
	Fail(name string)
	Success(name string)

	Duration(name string, start time.Time)
}

type ProcedureMeter interface {
	ProcedureSuccess(procName string, start time.Time)
	ProcedureFail(procName string, start time.Time)
}

// EventMeter - this looks a lot like logging. do we want this?
type CacheMeter interface {
	Hit(name string)
	Miss(name string)
	Size(name string, size int)
	Duration(name string, start time.Time)
}

type Metrics struct {
	//Name - the name of this group of metrics
	Name        string
	Tags        []string
	Hostname    string
	Environment string
}

//Basic Meter

func (m *Metrics) Occurence(name string) {
}

// Duration - Called after a process has completed, captures both that an
// occurence of a process happened and how long it took.
func (m *Metrics) Duration(name string, start time.Time) {
}

// Distribution - Capture the distribution of a value. Mean, Max, Min, Avg, 95%
func (m *Metrics) Distribution(name string, value float64) {

}

// Current - A value where the actual value is the only thing that matters. A level or gauge.
func (m *Metrics) Value(name string, value float64) {

}

func (m *Metrics) Size(name string, size int64) {

}

///Cache Meter

// Hit - Cache hit
func (m *Metrics) Hit(name string) {
}

// Miss - Cache miss
func (m *Metrics) Miss(name string) {
}
