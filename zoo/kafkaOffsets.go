package zoo

type KafkaOffsets struct {
	ZookeeperHosts []string
	AppName        string
}

func (ko *KafkaOffsets) Offsets(topic string) (map[int]int64, error) {
	hosts := ko.ZookeeperHosts
	tc := &TopicConsumer{
		Topic: topic,
		App:   ko.AppName,
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
		Topic: topic,
		App:   ko.AppName,
	}
	ks := &KafkaState{
		Topic: topic,
	}
	keeper := NewKafkaKeeper(hosts, tc, ks)
	return keeper.SetOffset(partition, offset)
}
