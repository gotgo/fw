package redisc

type ZAddArgs struct {
	Key     string
	Members []*ScoredMember
}

type ZRevRangeArgs struct {
	Key         string
	Start, Stop int
}

type ZIncrByArgs struct {
	Key    string
	Amount int
	Member []byte
}

type ZCardArgs struct {
	Key string
}

type RedisSortedSetMock struct {
	ZAddArgs      []*ZAddArgs
	ZRevRangeArgs []*ZRevRangeArgs
	ZIncrByArgs   []*ZIncrByArgs
	ZCardArgs     []*ZCardArgs
}

//TODO: enable callers to control return values

func NewRedisSortedSetMock() *RedisSortedSetMock {
	return &RedisSortedSetMock{
		ZAddArgs:      make([]*ZAddArgs, 0),
		ZRevRangeArgs: make([]*ZRevRangeArgs, 0),
		ZIncrByArgs:   make([]*ZIncrByArgs, 0),
		ZCardArgs:     make([]*ZCardArgs, 0),
	}
}

func (m *RedisSortedSetMock) ZAdd(key string, members []*ScoredMember) (int, error) {
	m.ZAddArgs = append(m.ZAddArgs, &ZAddArgs{key, members})
	return 1, nil
}

// ZRevRange returns a subset ordered in descending order
func (m *RedisSortedSetMock) ZRevRange(key string, start, stop int) ([]*ScoredMember, error) {
	m.ZRevRangeArgs = append(m.ZRevRangeArgs, &ZRevRangeArgs{key, start, stop})
	return []*ScoredMember{}, nil
}

func (rc *RedisSortedSetMock) ZRevRangeByScore(key string, max, min int) ([]*ScoredMember, error) {
	//TODO
	//m.ZRevRangeArgs = append(m.ZRevRangeArgs, &ZRevRangeArgs{key, start, stop})
	return []*ScoredMember{}, nil
}
func (m *RedisSortedSetMock) ZIncrBy(key string, amount int, member []byte) (int, error) {
	m.ZIncrByArgs = append(m.ZIncrByArgs, &ZIncrByArgs{key, amount, member})
	return 1, nil
}

// ZCard returns the Cardinality (i.e. count) of the set
func (m *RedisSortedSetMock) ZCard(key string) (int, error) {
	m.ZCardArgs = append(m.ZCardArgs, &ZCardArgs{key})
	return 1, nil
}
