package redisc

import (
	"encoding/json"
	"errors"
	"math/rand"
	"time"

	"github.com/amattn/deeperror"
	"github.com/garyburd/redigo/redis"
	"github.com/gotgo/fw/logging"
)

type RedisCache struct {
	readPool  *redis.Pool
	writePool *redis.Pool
	Log       logging.Logger `inject:""`
	Encoder   func(v interface{}) ([]byte, error)
	Decoder   func(data []byte, v interface{}) error
}

// NewRedisCache creates a new cache service connecting to the given
// hostUris and hostPassword. If there is no hostPassword, then pass an empty string.
func NewService(readUris []string, writeUris []string, hostPassword string) (*RedisCache, error) {
	r := new(RedisCache)
	r.readPool = r.newPool(readUris, hostPassword)
	r.writePool = r.newPool(writeUris, hostPassword)
	if err := r.Ping(); err != nil {
		return nil, err
	} else {
		return r, nil
	}
}

func (r *RedisCache) Ping() error {
	if conn, err := r.read(); err != nil {
		return err
	} else {
		defer conn.Close()
		if reply, err := redis.String(conn.Do("PING")); err != nil {
			return err
		} else if reply == "PONG" {
			return nil
		} else {
			return errors.New("unexpected reply " + reply)
		}
	}
}

func (r *RedisCache) GetBytes(ns, key string) (result []byte, err error) {
	if conn, err := r.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()

		useKey := getKey(ns, key)
		reply, err := conn.Do("GET", useKey)
		if err != nil {
			return nil, err
		} else if reply == nil {
			//miss
			return nil, nil
		}

		if bytes, err := redis.Bytes(reply, nil); err != nil {
			return nil, err
		} else {
			return bytes, nil
		}
	}
}

func (r *RedisCache) MGet(ns string, keys []string) (result []string, err error) {
	conn, err := r.read()
	if err != nil {
		return nil, deeperror.New(rand.Int63(), "Redis connect fail", err)
	}
	defer conn.Close()

	useKeys := getKeys(ns, keys)
	if values, err := arrayOfStrings(conn.Do("MGET", useKeys...)); err != nil {
		return nil, deeperror.New(rand.Int63(), "MGET fail", err)
	} else {
		return values, nil
	}
}

type KeyValueString struct {
	Key   string
	Value string
}

func flatten(ns string, kvs []*KeyValueString) []interface{} {
	r := make([]interface{}, 2*len(kvs))
	for i, kv := range kvs {
		j := i * 2
		r[j] = getKey(ns, kv.Key)
		r[j+1] = kv.Value
	}
	return r
}

func (r *RedisCache) MSet(ns string, kv []*KeyValueString) error {
	conn, err := r.write()
	if err != nil {
		return deeperror.New(rand.Int63(), "Redis connect fail", err)
	}
	defer conn.Close()

	if _, err := conn.Do("MSET", flatten(ns, kv)...); err != nil {
		return deeperror.New(rand.Int63(), "MSET fail", err)
	}
	return nil
}

func (r *RedisCache) SetBytes(ns, key string, bytes []byte) error {
	if conn, err := r.write(); err != nil {
		return err
	} else {
		defer conn.Close()
		useKey := getKey(ns, key)

		if _, err = redis.String(conn.Do("SET", useKey, bytes)); err != nil {
			return err
		}
		return nil
	}
}

// Get value from cache by key.
func (r *RedisCache) Get(ns, key string, instance interface{}) (miss bool, err error) {
	if bytes, err := r.GetBytes(ns, key); err != nil {
		return true, err
	} else if bytes == nil {
		return true, nil
	} else if err = r.unmarshal(bytes, &instance); err != nil {
		return true, err
	}
	return false, nil
}

// Set a value in cache by the given key.
func (r *RedisCache) Set(ns, key string, instance interface{}) error {
	if conn, err := r.write(); err != nil {
		return err
	} else {
		defer conn.Close()

		if value, err := r.marshal(instance); err != nil {
			return err
		} else if _, err = redis.String(conn.Do("SET", key, value)); err != nil {
			return err
		}
		return nil
	}
}

func (r *RedisCache) Increment(hashName, fieldName string, by int) (int64, error) {
	if conn, err := r.write(); err != nil {
		return -1, err
	} else {
		defer conn.Close()
		return redis.Int64(conn.Do("HINCRBY", hashName, fieldName, by))
	}
}

func (r *RedisCache) GetHashInt64(hashName string) (map[string]int64, error) {
	if conn, err := r.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		if list, err := redis.Values(conn.Do("HINCRBY", hashName)); err != nil {
			return nil, err
		} else {
			hash := make(map[string]int64)
			for i := 0; i < len(list); i += 2 {
				key, _ := redis.String(list[i], nil)
				hash[key], _ = redis.Int64(list[i+1], nil)
			}
			return hash, nil
		}
	}
}

func (r *RedisCache) write() (redis.Conn, error) {
	return r.connection(true)
}
func (r *RedisCache) read() (redis.Conn, error) {
	return r.connection(false)
}
func (r *RedisCache) connection(write bool) (redis.Conn, error) {
	pool := r.readPool
	if write {
		pool = r.writePool
	}
	if pool == nil {
		return nil, errors.New("pool is null")
	}
	conn, err := pool.Dial()
	if err != nil {
		r.Log.Error("Redis Dial Error", err)
		return nil, err
	}
	return conn, nil
}

func arrayOfBytes(results interface{}, err error) ([][]byte, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, err
	} else {
		result := make([][]byte, len(values))
		for i := 0; i < len(values); i++ {
			result[i] = values[i].([]byte)
		}
		return result, nil
	}
}

func arrayOfStrings(results interface{}, err error) ([]string, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, deeperror.New(rand.Int63(), "redis values fail", err)
	} else {
		result := make([]string, len(values))
		for i, value := range values {
			result[i], _ = redis.String(value, nil)
		}
		return result, nil
	}
}

func (r *RedisCache) marshal(v interface{}) ([]byte, error) {
	encoder := r.Encoder
	if encoder == nil {
		encoder = json.Marshal
	}
	return encoder(v)
}

func (r *RedisCache) unmarshal(data []byte, v interface{}) error {
	decoder := r.Decoder
	if decoder == nil {
		decoder = json.Unmarshal
	}
	return decoder(data, &v)
}

func (r *RedisCache) newPool(servers []string, password string) *redis.Pool {
	pickServer := func() string {
		i := rand.Int31n(int32(len(servers)))
		return servers[i]
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			server := pickServer()
			c, err := redis.Dial("tcp", server)
			if err != nil {
				r.Log.Error("can not dial redis instance: "+server, err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					r.Log.Error("incorrect redis password for instance: "+server, err)
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
