package multi

import "sync"

func NewDataContext() *DataContext {
	return &DataContext{
		data: make(map[string]interface{}),
	}
}

type DataContext struct {
	data  map[string]interface{}
	mutex sync.Mutex
}

func (c *DataContext) Set(key string, value interface{}) {

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.data[key] = value
}

func (c *DataContext) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	r := c.data[key]
	return r
}
