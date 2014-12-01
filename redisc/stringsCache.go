package redisc

import (
	"math/rand"

	"github.com/amattn/deeperror"
	"github.com/gotgo/fw/me"
)

func (r *RedisCache) MGet(keys []string) (result []string, err error) {
	conn, err := r.read()
	if err != nil {
		return nil, me.Err(err, "Redis connect fail")
	}
	defer conn.Close()

	useKeys := stringsToInterfaces(keys)
	if values, err := arrayOfStrings(conn.Do("MGET", useKeys...)); err != nil {
		return nil, deeperror.New(rand.Int63(), "MGET fail", err)
	} else {
		return values, nil
	}
}
func (r *RedisCache) MSet(ns string, kv []*KeyValueString) error {
	conn, err := r.write()
	if err != nil {
		return me.Err(err, "Redis connect fail")
	}
	defer conn.Close()

	if _, err := conn.Do("MSET", flatten(kv)...); err != nil {
		return deeperror.New(rand.Int63(), "MSET fail", err)
	}
	return nil
}

// Get value from cache by key.
func (r *RedisCache) Get(ns, key string, instance interface{}) (miss bool, err error) {
	if bytes, err := r.GetBytes(key); err != nil {
		return true, err
	} else if bytes == nil {
		return true, nil
	} else if err = r.unmarshal(bytes, &instance); err != nil {
		return true, err
	}
	return false, nil
}

func (r *RedisCache) SetNX(key string, value string) error {
	return r.setWithOverwrite(key, value, false)
}

func (r *RedisCache) Set(key string, value string) error {
	return r.setWithOverwrite(key, value, true)
}

// Set a value in cache by the given key.
func (r *RedisCache) setWithOverwrite(key string, value string, overwrite bool) error {
	if conn, err := r.write(); err != nil {
		return err
	} else {
		defer conn.Close()

		command := "SET"
		if overwrite == false {
			command = "SETNX"
		}
		if _, err = conn.Do(command, key, value); err != nil {
			return err
		}
		return nil
	}
}
