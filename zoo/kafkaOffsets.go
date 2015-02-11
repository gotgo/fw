package zoo

type KafkaOffsets struct {
	ZookeeperHosts []string
	// ConsumerApp - consumer application get's it's own distinct offsets on the topic
	ConsumerApp string
}

func (ko *KafkaOffsets) Offsets(topic string) ([]*PartitionOffset, error) {
	hosts := ko.ZookeeperHosts
	tc := &TopicConsumer{
		Topic:       topic,
		ConsumerApp: ko.ConsumerApp,
	}
	ks := &KafkaState{
		Topic: topic,
	}
	keeper := NewKafkaKeeper(hosts, tc, ks)
	return keeper.GetOffsets()
}

func (ko *KafkaOffsets) SetOffset(topic string, partition int32, offset int64) error {
	hosts := ko.ZookeeperHosts
	tc := &TopicConsumer{
		Topic:       topic,
		ConsumerApp: ko.ConsumerApp,
	}
	ks := &KafkaState{
		Topic: topic,
	}
	keeper := NewKafkaKeeper(hosts, tc, ks)
	return keeper.SetOffset(partition, offset)
}
