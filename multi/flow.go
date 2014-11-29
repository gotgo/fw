package multi

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

func (f *Flow) NewStep(name string, action StepAction, state interface{}) {
	f.Steps[name] = &StepState{
		Name:   name,
		Action: action,
		State:  state,
	}
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
