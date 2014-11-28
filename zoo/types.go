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
	App   string
	// Root - suggested name /kafka-topics
	Root string
}

func (tc *TopicConsumer) basePath() string {
	return path.Join("/", tc.Root, tc.Topic, "apps", tc.App)
}

func (tc *TopicConsumer) PartitionsPath() string {
	// /{root}/{topic}/apps/{app}/partitions
	return path.Join(tc.basePath(), "partitions")
}
