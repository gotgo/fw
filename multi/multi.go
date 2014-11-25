package multi

type StepAction interface {
	Action(tuple DataFlow, done func())
	Error() error
}

func NewFlow() *Flow {
	return &Flow{
		Steps: make(map[string]*StepState),
		Data: &data{
			data: make(map[string]interface{}),
		},
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
	Data  DataFlow
}

func (f *Flow) StepsWithErrors() []*StepState {
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
	}
}

type StepState struct {
	Name     string
	Attempts int
	Action   StepAction
	State    interface{}
}

type DataFlow interface {
	Set(name string, value interface{})
	Get(name string) interface{}
}

type data struct {
	data map[string]interface{}
}

func (d *data) Set(name string, value interface{}) {
	d.data[name] = value
}
func (d *data) Get(name string) interface{} {
	return d.data[name]
}

func GatherFailures(cs ...*Coordinator) []*Flow {
	flows := []*Flow{}
	for _, c := range cs {
		for f := range c.Fail {
			flows = append(flows, f)
		}
	}
	return flows
}
