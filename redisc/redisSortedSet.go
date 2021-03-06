package redisc

import (
	"math"
	"strconv"

	"github.com/garyburd/redigo/redis"
)

func (rc *RedisCache) ZAdd(key string, members ...*ScoredMember) (int, error) {
	if len(members) == 0 {
		return 0, nil
	}

	if conn, err := rc.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()

		count := (len(members) * 2) + 1
		command := make([]interface{}, count)
		command[0] = key
		items := command[1:]

		for i := 0; i < len(members); i++ {
			j := i * 2
			items[j] = members[i].Score
			items[j+1] = members[i].Member
		}

		if added, err := redis.Int(conn.Do("ZADD", command...)); err != nil {
			rc.Log.Error("Redis ZADD fail", err)
			return 0, err
		} else {
			return added, nil
		}
	}
}

// ZRevRange returns a subset ordered in descending order
func (rc *RedisCache) ZRevRange(key string, start, stop int) ([]*ScoredMember, error) {
	if conn, err := rc.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		return rc.scoredMembers(conn.Do("ZREVRANGE", key, start, stop, "WITHSCORES"))
	}
}
func (rc *RedisCache) ZScore(key, member string) (*int, error) {
	if conn, err := rc.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		v, err := redis.Int(conn.Do("ZSCORE", key, member))
		if err == redis.ErrNil {
			return nil, nil
		} else if err != nil {
			return nil, err
		} else {
			return &v, err
		}
	}
}

func (rc *RedisCache) ZRevRangeByScore(key string, max, min int) ([]*ScoredMember, error) {
	if conn, err := rc.read(); err != nil {
		return nil, err
	} else {
		defer conn.Close()
		maxVal := strconv.Itoa(max)
		if max == math.MaxInt32 {
			maxVal = "+inf"
		}
		minVal := strconv.Itoa(min)
		if min == math.MinInt32 {
			minVal = "-inf"
		}
		return rc.scoredMembers(conn.Do("ZREVRANGEBYSCORE", key, maxVal, minVal, "WITHSCORES"))
	}
}
func (rc *RedisCache) ZIncrBy(key string, amount int, member string) (int, error) {
	if conn, err := rc.write(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("ZINCRBY", key, amount, member))
	}
}

// ZCard returns the Cardinality (i.e. count) of the set
func (rc *RedisCache) ZCard(key string) (int, error) {
	if conn, err := rc.read(); err != nil {
		return 0, err
	} else {
		defer conn.Close()
		return redis.Int(conn.Do("ZCard", key))
	}
}

func (rc *RedisCache) scoredMembers(results interface{}, err error) ([]*ScoredMember, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, err
	} else {
		result := make([]*ScoredMember, len(values)/2)
		for i := range result {
			j := i * 2

			mv, _ := redis.String(values[j], nil)
			score, _ := redis.Int(values[j+1], nil)
			member := &ScoredMember{
				Member: mv,
				Score:  score,
			}
			result[i] = member
		}
		return result, nil
	}
}
