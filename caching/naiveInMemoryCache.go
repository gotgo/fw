package caching

import (
	"encoding/json"

	"github.com/gotgo/fw/logging"
)

// NaiveInMemoryCache is used for TESTING! It has no TTL and will eventually eat all memory
type NaiveInMemoryCache struct {
	data map[string][]byte
	Log  logging.Logger
}

func NewNaiveInMemoryCache() *NaiveInMemoryCache {
	cache := new(NaiveInMemoryCache)
	cache.data = make(map[string][]byte)
	cache.Log = new(logging.NoOpLogger)
	return cache
}

func (simc *NaiveInMemoryCache) Get(ns, key string, v interface{}) (miss bool, err error) {
	useKey := getKey(ns, key)
	bytes := simc.data[useKey]

	if bytes == nil {
		simc.Log.Debugf("Cache MISS for Key='%s'", useKey)
		return true, nil
	} else {
		json.Unmarshal(bytes, &v)
		simc.Log.Debugf("Cache HIT for Key='%s' Value='%v'", useKey, v)
		return false, nil
	}
}

func (simc *NaiveInMemoryCache) Set(ns, key string, v interface{}) error {
	useKey := getKey(ns, key)

	if bytes, err := json.Marshal(v); err != nil {
		return err
	} else {
		simc.data[useKey] = bytes
		simc.Log.Debugf("Cache Set for Key='%s'", useKey)
		return nil
	}
}
