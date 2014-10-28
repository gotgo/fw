package main

type Manager struct {
	configs map[string]interface{}
	updates map[string]func(interface{})
}

func NewManager() *Manager {
	m := &Manager{
		configs: make(map[string]interface{}),
		updates: make(map[string]func(interface{})),
	}
	return m
}

func (m *Manager) Use(key string, cfg interface{}) {
	c := m.configs[key]
	m.configs[key] = cfg
	if c != nil {
		f := m.updates[key]
		if f != nil {
			f(cfg)
		}
	}
}
func (m *Manager) Get(key string) interface{} {
	return m.configs[key]
}

func (m *Manager) NotifyMe(key string, method func(newCfg interface{})) {
	m.updates[key] = method
}
