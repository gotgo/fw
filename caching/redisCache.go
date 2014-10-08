package caching

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gotgo/fw/logging"
)

type RedisCache struct {
	Pool    *redis.Pool
	Log     logging.Logger `inject:""`
	Encoder func(v interface{}) ([]byte, error)
	Decoder func(data []byte, v interface{}) error
}

// NewRedisCache creates a new cache service connecting to the given
// hostUri and hostPassword. If there is no hostPassword, then pass an empty string.
func NewService(hostUri, hostPassword string) *RedisCache {
	s := new(RedisCache)
	s.newPool(hostUri, hostPassword)
	return s
}

func getKey(ns, key string) string {
	useKey := fmt.Sprintf("%s:%s", ns, key)
	return useKey
}

func (s *RedisCache) GetBytes(ns, key string) (result []byte, err error) {
	if conn, err := s.connection(); err != nil {
		return nil, err
	} else {
		defer conn.Close()

		useKey := getKey(ns, key)
		reply, err := conn.Do("GET", useKey)
		if err != nil {
			s.Log.Error("Redis GET failed", err)
			return nil, err
		} else if reply == nil {
			//miss
			return nil, nil
		}

		if bytes, err := redis.Bytes(reply, nil); err != nil {
			s.Log.Error("Redis GET failed", err)
			return nil, err
		} else {
			return bytes, nil
		}
	}
}

func (s *RedisCache) SetBytes(ns, key string, bytes []byte) error {
	if conn, err := s.connection(); err != nil {
		return err
	} else {
		defer conn.Close()
		useKey := getKey(ns, key)

		if _, err = redis.String(conn.Do("SET", useKey, bytes)); err != nil {
			s.Log.Error("Redis SET failed", err)
			return err
		}
		return nil
	}
}

// Get value from cache by key.
func (cs *RedisCache) Get(ns, key string, instance interface{}) (miss bool, err error) {
	if bytes, err := cs.GetBytes(ns, key); err != nil {
		return true, err
	} else if bytes == nil {
		return true, nil
	} else if err = cs.unmarshal(bytes, &instance); err != nil {
		cs.Log.UnmarshalFail("Redis GET unmarshal fail", err)
		return true, err
	}
	return false, nil
}

// Set a value in cache by the given key.
func (s *RedisCache) Set(ns, key string, instance interface{}) error {
	if conn, err := s.connection(); err != nil {
		return err
	} else {
		defer conn.Close()

		if value, err := s.marshal(instance); err != nil {
			s.Log.MarshalFail("Redis marshal error", err)
			return err
		} else if _, err = redis.String(conn.Do("SET", key, value)); err != nil {
			s.Log.Error("Redis SET fail", err)
			return err
		}
		return nil
	}
}

func (s *RedisCache) SetAdd(listKey string, member []byte) (int, error) {
	if conn, err := s.connection(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("SADD", listKey, member))
	}
}

func (s *RedisCache) SetMembers(listKey string) ([][]byte, error) {
	if conn, err := s.connection(); err != nil {
		return [][]byte{}, err
	} else {
		defer conn.Close()
		if values, err := redis.Values(conn.Do("SMEMBERS", listKey)); err != nil {
			return nil, err
		} else {
			result := make([][]byte, len(values))
			for i := 0; i < len(values); i++ {
				result[i] = values[i].([]byte)
			}
			return result, nil
		}
	}
}

func (s *RedisCache) SetRemove(listKey string, member []byte) (int, error) {
	if conn, err := s.connection(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("SREM", listKey, member))
	}
}

func (s *RedisCache) Increment(hashName, fieldName string, by int) (int64, error) {
	if conn, err := s.connection(); err != nil {
		return -1, err
	} else {
		defer conn.Close()
		return redis.Int64(conn.Do("HINCRBY", hashName, fieldName, by))
	}
}

func (s *RedisCache) GetHashInt64(hashName string) (map[string]int64, error) {
	if conn, err := s.connection(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		if list, err := redis.Values(conn.Do("HINCRBY", hashName)); err != nil {
			return nil, err
		} else {
			hash := make(map[string]int64)
			for i := 0; i < len(list); i += 2 {
				hash[(list[i]).(string)] = list[i+1].(int64)
			}
			return hash, nil
		}
	}
}

func (s *RedisCache) connection() (redis.Conn, error) {
	if s.Pool == nil {
		return nil, errors.New("pool is null")
	}
	conn, err := s.Pool.Dial()
	if err != nil {
		s.Log.Error("Redis Dial Error", err)
		return nil, err
	}
	return conn, nil
}
func (cs *RedisCache) marshal(v interface{}) ([]byte, error) {
	encoder := cs.Encoder
	if encoder == nil {
		encoder = json.Marshal
	}
	return encoder(v)
}

func (cs *RedisCache) unmarshal(data []byte, v interface{}) error {
	decoder := cs.Decoder
	if decoder == nil {
		decoder = json.Unmarshal
	}
	return decoder(data, &v)
}

func (cs *RedisCache) newPool(server, password string) {
	cs.Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				cs.Log.Error("can not dial redis instance: "+server, err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					cs.Log.Error("incorrect redis password for instance: "+server, err)
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			if err != nil {
				cs.Log.Warn("failed to ping redis instance", "instance", server)
			}
			return err
		},
	}
}
