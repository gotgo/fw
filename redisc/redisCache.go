package redisc

import (
	"encoding/json"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/amattn/deeperror"
	"github.com/garyburd/redigo/redis"
	"github.com/gotgo/fw/logging"
)

type RedisCache struct {
	readPool  *roundRobinPools
	writePool *roundRobinPools
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

func (r *RedisCache) GetBytes(key string) (result []byte, err error) {
	if conn, err := r.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()

		reply, err := conn.Do("GET", key)
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

type KeyValueString struct {
	Key   string
	Value string
}

func flatten(kvs []*KeyValueString) []interface{} {
	r := make([]interface{}, 2*len(kvs))
	for i, kv := range kvs {
		j := i * 2
		r[j] = kv.Key
		r[j+1] = kv.Value
	}
	return r
}

func (r *RedisCache) SetBytes(key string, bytes []byte) error {
	if conn, err := r.write(); err != nil {
		return err
	} else {
		defer conn.Close()
		if _, err = redis.String(conn.Do("SET", key, bytes)); err != nil {
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
	pools := r.readPool
	if write {
		pools = r.writePool
	}
	if pools == nil {
		return nil, errors.New("pool is null")
	}
	conn, err := pools.GetPool().Dial()
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

// as per @GaryBurd suggestion
type roundRobinPools struct {
	mu    sync.Mutex
	i     int
	pools []*redis.Pool
}

func (p *roundRobinPools) GetPool() *redis.Pool {
	if len(p.pools) == 0 {
		return nil
	}
	p.mu.Lock()
	p.i = (p.i + 1) % len(p.pools)
	p.mu.Unlock()
	return p.pools[p.i]
}

// creates a pool of connection pools
func (r *RedisCache) newPool(servers []string, password string) *roundRobinPools {
	pools := make([]*redis.Pool, len(servers))
	for i, s := range servers {
		pools[i] = &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				server := s
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
	return &roundRobinPools{
		pools: pools,
	}
}
