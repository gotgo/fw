package multi

type StepSource interface {
	Action(done func())
	Error() error
}

func NewFlow() *Flow {
	return &Flow{
		Steps: make(map[string]*StepState),
	}
}

func CreateFlows(count int) []*Flow {
	flows := make([]*Flow, count)
	for i := 0; i < count; i++ {
		flows[i] = NewFlow()
	}
	return flows
}

type Flow struct {
	Steps map[string]*StepState
}

func (f *Flow) WithErrors() []*StepState {
	steps := []*StepState{}
	for _, s := range f.Steps {
		if err := s.Source.Error(); err != nil {
			steps = append(steps, s)
		}
	}
	return steps
}

func (f *Flow) NewStep(name string, source StepSource, state interface{}) { //helper method
	f.Steps[name] = &StepState{
		Name:   name,
		Source: source,
		State:  state,
	}
}

type StepState struct {
	Name     string
	Attempts int
	Source   StepSource
	State    interface{}
}
