package stats

func Count(name string) {}
func Hit(name string)   {}
func Miss(name string)  {}

type Metrics struct {
}

func (m *Metrics) Count(name string) {

}

func (m *Metrics) Hit(name string) {
	//count
	//rate
	//%
}

func (m *Metrics) Miss(name string) {
	//count
	//rate
	//%
}
