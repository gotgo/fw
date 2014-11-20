package caching

import "encoding/json"

// NaiveInMemoryCache is used for TESTING! It has no TTL and will eventually eat all memory
type NaiveInMemoryCache struct {
	data map[string][]byte
}

func NewNaiveInMemoryCache() *NaiveInMemoryCache {
	cache := new(NaiveInMemoryCache)
	cache.data = make(map[string][]byte)
	return cache
}

func (simc *NaiveInMemoryCache) Get(ns, key string, v interface{}) (miss bool, err error) {
	useKey := getKey(ns, key)
	bytes := simc.data[useKey]

	if bytes == nil {
		return true, nil
	} else {
		json.Unmarshal(bytes, &v)
		return false, nil
	}
}

func (simc *NaiveInMemoryCache) Set(ns, key string, v interface{}) error {
	useKey := getKey(ns, key)

	if bytes, err := json.Marshal(v); err != nil {
		return err
	} else {
		simc.data[useKey] = bytes
		return nil
	}
}
