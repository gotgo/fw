package caching

// NoOpCache is used to disable caching. All operations miss or do nothing
type NoOpCache struct{}

func (noc *NoOpCache) Get(ns, key string, instance interface{}) (miss bool, err error) {
	return true, nil
}

func (noc *NoOpCache) Set(ns, key string, instance interface{}) error {
	return nil
}
