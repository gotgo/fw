package caching

type BytesCache interface {
	GetBytes(ns, key string) (result []byte, err error)
	SetBytes(ns, key string, bytes []byte) error
}

type ObjectCache interface {
	Get(ns, key string, instance interface{}) (miss bool, err error)
	Set(ns, key string, instance interface{}) error
}

type HashCache interface {
	Increment(hashName, fieldName string, by int) (int64, error)
	GetHashInt64(hashName string) (map[string]int64, error)
}

type ListCache interface {
	ListPushRight(listKey string, value []byte) (int64, error)

	ListPopLeft(listKey string) ([]byte, error)

	ListGetRange(listKey string, startIndex, endIndex int) ([]byte, error)
}

type SetCache interface {
	SetAdd(listKey string, member []byte) (int, error)
	SetMembers(listKey string) ([][]byte, error)
	SetRemove(listKey string, member []byte) (int, error)
}

type RedisSet interface {
	// Add
	SAdd(key string, items [][]byte) (int, error)
	// Remove
	SRem(key string, items [][]byte) (int, error)

	SRandMember(key string, count int) ([][]byte, error)

	SMembers(listKey string) ([][]byte, error)
}

type RedisSortedSet interface {
	ZAdd(key string, members []*ScoredMember)

	// ZRevRange returns a subset ordered in descending order
	ZRevRange(key string, start, stop int) ([]*ScoredMember, error)

	ZIncrBy(key string, amount int, member []byte) (int, error)

	// ZCard returns the Cardinality (i.e. count) of the set
	ZCard(key string) (int, error)
}
