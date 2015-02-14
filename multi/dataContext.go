package multi

import "sync"

func NewDataContext() *DataContext {
	return &DataContext{
		data:  make(map[string]interface{}),
		mutex: new(sync.Mutex),
	}
}

type DataContext struct {
	data  map[string]interface{}
	mutex *sync.Mutex
}

func (c *DataContext) Set(key string, value interface{}) {
	c.mutex.Lock()
	c.data[key] = value
	c.mutex.Unlock()
}

func (c *DataContext) Get(key string) interface{} {
	c.mutex.Lock()
	r := c.data[key]
	c.mutex.Unlock()
	return r
}
