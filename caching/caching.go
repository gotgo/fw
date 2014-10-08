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
