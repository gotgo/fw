package zoo

import "path"

type KafkaState struct {
	Topic string
	// Root - typically blank or "/kafka"
	Root string
}

func (ks *KafkaState) PartitionsPath() string {
	// /kafka/brokers/topics/{topic}/partitions/
	return path.Join("/", ks.Root, "brokers", "topics", ks.Topic, "partitions")
}

type TopicConsumer struct {
	Topic string
	// ConsumerApp - typically the name of the app or function for it's own dedicated set of offsets
	ConsumerApp string
	// Root - suggested name /kafka-topics
	Root string
}

func (tc *TopicConsumer) basePath() string {
	return path.Join("/", tc.Root, tc.Topic, "consumers", tc.ConsumerApp)
}

func (tc *TopicConsumer) PartitionsPath() string {
	// /{root}/{topic}/consumers/{consumer}/partitions
	return path.Join(tc.basePath(), "partitions")
}
