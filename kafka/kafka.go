package kafka

type Producer interface {
	Close()
	SendBytes(message []byte, topic, key string) error
}

type ConsumerEvent struct {
	Error     error
	Message   []byte
	Offset    int64
	Partition int32
}

type ConsumerChannel interface {
	Events() <-chan *ConsumerEvent
	Close() error
}

type Sender interface {
	SendBytes(message []byte, topic, key string) error
}

// DividePartitions - splits the number of partitions in to buckets numbering the splitBy
func DividePartitions(partitions, splitBy int) [][]int {
	buckets := make([][]int, splitBy, splitBy)

	if splitBy > partitions {
		splitBy = partitions
	}

	count := partitions / splitBy
	extra := partitions % splitBy

	for j := 0; j < splitBy; j++ {
		index := 0 + j
		bucketSize := count
		if j < extra {
			bucketSize += 1
		}

		bucket := make([]int, bucketSize, bucketSize)
		buckets[j] = bucket

		for i := 0; i < bucketSize; i++ {
			bucket[i] = index
			index += splitBy
		}
	}
	return buckets
}
