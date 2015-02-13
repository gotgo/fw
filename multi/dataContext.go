package multi

import "sync"

type DataContext struct {
	data  map[string]interface{}
	mutex sync.Mutex
}

// must be called from inside a lock
func (c *DataContext) unsafeGetData() map[string]interface{} {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	return c.data
}

func (c *DataContext) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data := c.unsafeGetData()
	data[key] = value
}

func (c *DataContext) Get(key string) interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	data := c.unsafeGetData()
	return data[key]
}
