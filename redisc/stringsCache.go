package redisc

import (
	"math/rand"
	"time"

	"github.com/amattn/deeperror"
	"github.com/garyburd/redigo/redis"
	"github.com/gotgo/fw/me"
)

func (r *RedisCache) MGet(keys []string) (result []string, err error) {
	conn, err := r.read()
	if err != nil {
		return nil, me.Err(err, "Redis connect fail")
	}
	defer conn.Close()

	useKeys := StringsToInterfaces(keys)
	if values, err := ArrayOfStrings(conn.Do("MGET", useKeys...)); err != nil {
		return nil, deeperror.New(rand.Int63(), "MGET fail", err)
	} else {
		return values, nil
	}
}

func (r *RedisCache) MSet(kv []*KeyValueString) error {
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
func (r *RedisCache) Get(key string) (value string, err error) {
	if conn, err := r.read(); err != nil {
		return "", err
	} else {
		defer conn.Close()
		return redis.String(conn.Do("GET", key))
	}
}

func (r *RedisCache) SetNX(key string, value string) error {
	return r.setWithOverwrite(key, value, false)
}

func (r *RedisCache) Set(key string, value string) error {
	return r.setWithOverwrite(key, value, true)
}

type SetP struct {
	Key   string
	Value string
	TTL   time.Duration
	NX    bool
	XX    bool
}

func (s *SetP) Command() (string, []interface{}) {
	v := []interface{}{s.Key, s.Value}
	if s.TTL > 0 {
		v = append(v, "EX")
		v = append(v, s.TTL*time.Second)
	}
	return "SET", v
}

/////////////////////////
// NEW INTERFACE??

func (r *RedisCache) Write(command string, args ...interface{}) (interface{}, error) {
	//if len(args) == 1 && args[0]
	if conn, err := r.write(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		return conn.Do(command, args...)
	}
}

func (r *RedisCache) Read(command string, args ...interface{}) (interface{}, error) {
	if conn, err := r.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		return conn.Do(command, args...)
	}
}
func (r *RedisCache) ReadInt64(command string, args ...interface{}) (int64, error) {
	if conn, err := r.read(); err != nil {
		return -1, err
	} else {
		defer conn.Close()
		return redis.Int64(conn.Do(command, args...))
	}
}

//////////////////////////////////

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
