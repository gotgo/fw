package multi

type StepAction interface {
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
		if err := s.Action.Error(); err != nil {
			steps = append(steps, s)
		}
	}
	return steps
}

func (f *Flow) NewStep(name string, action StepAction, state interface{}) { //helper method
	f.Steps[name] = &StepState{
		Name:   name,
		Action: action,
		State:  state,
		Output: make(map[string]interface{}),
	}
}

type StepState struct {
	Name     string
	Attempts int
	Action   StepAction
	State    interface{}
	Output   map[string]interface{}
}

type DataFlow interface {
	Output(name string, value interface{})
	Input(name string) interface{}
}
