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

type data struct {
	data map[string]interface{}
}

func (d *data) Set(name string, value interface{}) {
	d.data[name] = value
}

func (d *data) Get(name string) interface{} {
	return d.data[name]
}
