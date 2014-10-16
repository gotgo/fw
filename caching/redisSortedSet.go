package caching

import "github.com/garyburd/redigo/redis"

type ScoredMember struct {
	Score  int
	Member []byte
}

func (rc *RedisCache) ZAdd(key string, members []*ScoredMember) (int, error) {
	if conn, err := rc.connection(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		items := make([]interface{}, len(members)*2)
		for i := range members {
			j := i * 2
			items[j] = members[i].Score
			items[j+1] = members[i].Member
		}

		if added, err := redis.Int(conn.Do("ZADD", key, items)); err != nil {
			rc.Log.Error("Redis ZADD fail", err)
			return 0, err
		} else {
			return added, nil
		}
	}
}

// ZRevRange returns a subset ordered in descending order
func (rc *RedisCache) ZRevRange(key string, start, stop int) ([]*ScoredMember, error) {
	if conn, err := rc.connection(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		return scoredMembers(conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	}
}

func (rc *RedisCache) ZIncrBy(key string, amount int, member []byte) (int, error) {
	if conn, err := rc.connection(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("ZINCRBY", key, amount, member))
	}
}

// ZCard returns the Cardinality (i.e. count) of the set
func (rc *RedisCache) ZCard(key string) (int, error) {
	if conn, err := rc.connection(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("ZCard", key))
	}
}

func scoredMembers(results interface{}, err error) ([]*ScoredMember, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, err
	} else {
		result := make([]*ScoredMember, len(values)/2)
		for i := range result {
			j := i * 2
			member := &ScoredMember{
				Member: values[j].([]byte),
				Score:  values[j+1].(int),
			}
			result[i] = member
		}
		return result, nil
	}
}
