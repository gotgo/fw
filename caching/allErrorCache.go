package caching

import "errors"

// AllErrorCache is used for testing.  Every interface that can return an error, does
type AllErrorCache struct{}

func (simc *AllErrorCache) Get(ns, key string, v interface{}) (miss bool, err error) {
	return true, errors.New("test error")
}

func (simc *AllErrorCache) Set(ns, key string, v interface{}) error {
	return errors.New("test error")
}
