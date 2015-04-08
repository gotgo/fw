package redisc

import (
	"math/rand"

	"github.com/amattn/deeperror"
	"github.com/garyburd/redigo/redis"
)

type StringsCache interface {
	Get(key string) (value string, err error)
	Set(key string, value string) error
	SetNX(key string, value string) error
	MGet(keys []string) (result []string, err error)
	MSet(kv []*KeyValueString) error
}

type Set interface {
	// Add
	SAdd(key string, items ...string) (int, error)
	// Remove
	SRem(key string, items ...string) (int, error)

	SRandMember(key string, count int) ([]string, error)

	SMembers(listKey string) ([]string, error)
}

type SortedSet interface {
	ZAdd(key string, members ...*ScoredMember) (int, error)

	// ZRevRange returns a subset ordered in descending order
	ZRevRange(key string, start, stop int) ([]*ScoredMember, error)
	ZRevRangeByScore(key string, max, min int) ([]*ScoredMember, error)

	ZIncrBy(key string, amount int, member string) (int, error)

	// ZCard returns the Cardinality (i.e. count) of the set
	ZCard(key string) (int, error)

	// ZSScore if the member or key does not exists return a nil int with no error
	ZScore(key, member string) (*int, error)
}

type Client interface {
	StringsCache
	SortedSet
	Set
	Write(command string, args ...interface{}) (interface{}, error)
	Read(command string, args ...interface{}) (interface{}, error)
}

type ScoredMember struct {
	Score  int
	Member string
}

func StringsToInterfaces(keys []string) []interface{} {
	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = k
	}
	return result
}

func PrefixStringsToInterfaces(prefix string, keys []string) []interface{} {
	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = prefix + k
	}
	return result
}

func GetMembers(members []*ScoredMember) []string {
	keys := make([]string, len(members))
	for i, m := range members {
		keys[i] = m.Member
	}
	return keys
}

func Prefix(namespace string, keys []string) []string {
	result := make([]string, len(keys))
	for i, k := range keys {
		result[i] = namespace + k
	}
	return result
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

func ArrayOfBytes(results interface{}, err error) ([][]byte, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, err
	} else {
		if values == nil {
			return [][]byte{}, nil
		} else {
			result := make([][]byte, len(values))
			for i := 0; i < len(values); i++ {
				if values[i] != nil {
					result[i] = values[i].([]byte)
				}
			}
			return result, nil
		}
	}
}

func ArrayOfStrings(results interface{}, err error) ([]string, error) {
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

func ArrayOfInts(results interface{}, err error) ([]int, error) {
	if values, err := redis.Values(results, err); err != nil {
		return nil, deeperror.New(rand.Int63(), "redis values fail", err)
	} else {
		result := make([]int, len(values))
		for i, value := range values {
			result[i], _ = redis.Int(value, nil)
		}
		return result, nil
	}
}
