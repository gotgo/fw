package config

import "reflect"

type Manager struct {
	configs map[reflect.Type]interface{}
}

func NewManager() *Manager {
	m := &Manager{
		configs: make(map[reflect.Type]interface{}),
	}
	return m
}

func (m *Manager) Use(instance interface{}) {
	if instance == nil {
		panic("manager instance can not be nil")
	}
	key := reflect.TypeOf(instance)
	m.configs[key] = instance
}

func (m *Manager) Get(instance interface{}) {
	key := reflect.TypeOf(instance)
	i := m.configs[key]
	instance = i
}

type ConfigProvider interface {
	Get(instance interface{})
}
