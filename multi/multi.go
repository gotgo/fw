package multi

type StepAction interface {
	Action(tuple DataFlow, done func())
	Error() error
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

func GatherFailures(cs ...*Coordinator) []*Flow {
	flows := []*Flow{}
	for _, c := range cs {
		for f := range c.Fail {
			flows = append(flows, f)
		}
	}
	return flows
}

type TaskAction interface {
	Run(input interface{}) (interface{}, error)
}
