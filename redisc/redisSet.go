package redisc

import (
	"github.com/garyburd/redigo/redis"
	"github.com/gotgo/fw/me"
)

// SAdd returns the number of items added to the set from the items given.
func (rc *RedisCache) SAdd(key string, items ...string) (int, error) {
	if conn, err := rc.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		if added, err := redis.Int(conn.Do("SADD", join(key, items))); err != nil {
			return 0, me.Err(err, "redis SADD fail")
		} else {
			return added, nil
		}
	}
}

func join(key string, items []string) []interface{} {
	result := make([]interface{}, len(items)+1)
	result[0] = key
	for i, item := range items {
		result[i+1] = item
	}
	return result
}

// SRem returns the number items removed from the set
func (rc *RedisCache) SRem(key string, items ...string) (int, error) {
	if conn, err := rc.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		if removed, err := redis.Int(conn.Do("SREM", join(key, items))); err != nil {
			return 0, me.Err(err, "Redis SREM fail")
		} else {
			return removed, nil
		}
	}
}

func (rc *RedisCache) SRandMember(key string, count int) ([]string, error) {
	if conn, err := rc.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		return arrayOfStrings(conn.Do("SRANDMEMBER", key))
	}
}

func (s *RedisCache) SMembers(listKey string) ([]string, error) {
	if conn, err := s.read(); err != nil {
		return []string{}, err
	} else {
		defer conn.Close()
		return arrayOfStrings(conn.Do("SMEMBERS", listKey))
	}
}

/////////////////////////////////////////////

// SAdd returns the number of items added to the set from the items given.
func (rc *RedisCache) SetAdd(key string, items []interface{}) (int, error) {
	toSend := make([][]byte, len(items))
	for i := range items {
		item := items[i]
		if value, err := rc.marshal(item); err != nil {
			rc.Log.MarshalFail("Redis SADD", item, err)
			return 0, err
		} else {
			toSend[i] = value
		}
	}
	if conn, err := rc.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		if added, err := redis.Int(conn.Do("SADD", key, toSend)); err != nil {
			rc.Log.Error("Redis SADD fail", err)
			return 0, err
		} else {
			return added, nil
		}
	}
}

func (s *RedisCache) SetRemove(listKey string, member []byte) (int, error) {
	if conn, err := s.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("SREM", listKey, member))
	}
}
