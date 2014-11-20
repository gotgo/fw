package redisc

import "fmt"

type StringsCache interface {
	MGet(ns string, keys []string) (result []string, err error)
}

type Set interface {
	// Add
	SAdd(key string, items [][]byte) (int, error)
	// Remove
	SRem(key string, items [][]byte) (int, error)

	SRandMember(key string, count int) ([][]byte, error)

	SMembers(listKey string) ([][]byte, error)
}

type SortedSet interface {
	ZAdd(key string, members []*ScoredMember) (int, error)

	// ZRevRange returns a subset ordered in descending order
	ZRevRange(key string, start, stop int) ([]*ScoredMember, error)

	ZIncrBy(key string, amount int, member []byte) (int, error)

	// ZCard returns the Cardinality (i.e. count) of the set
	ZCard(key string) (int, error)
}

type Client interface {
	StringsCache
	SortedSet
}

type ScoredMember struct {
	Score  int
	Member string
}

func getKey(ns, key string) string {
	useKey := fmt.Sprintf("%s:%s", ns, key)
	return useKey
}

func getKeys(ns string, keys []string) []interface{} {
	result := make([]interface{}, len(keys))
	for i, k := range keys {
		result[i] = getKey(ns, k)
	}
	return result
}
