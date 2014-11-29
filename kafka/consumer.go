package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/gotgo/fw/logging"
)

type Consumer struct {
	ClientName      string
	StartingOffsets map[int]int64
	Log             logging.Logger `inject:""`
	Topic           string
	stopper         chan struct{}
	events          chan *ConsumerEvent
	client          *sarama.Client
}

func NewConsumer(client *sarama.Client, clientName, topic string, startingOffsets map[int]int64) *Consumer {
	consumer := &Consumer{
		ClientName:      clientName,
		StartingOffsets: startingOffsets,
		Topic:           topic,
		events:          make(chan *ConsumerEvent),
		client:          client,
	}
	return consumer
}

// Run - Nonblocking
func (c *Consumer) Run() {
	consumers := c.setupConsumers()
	c.consume(consumers)
}

func (c *Consumer) Events() <-chan *ConsumerEvent {
	return c.events
}

func (c *Consumer) Shutdown() {
	if c.stopper != nil {
		close(c.stopper)
	}
}

func (c *Consumer) setupConsumers() []*sarama.Consumer {
	offsets := c.StartingOffsets
	count := len(offsets)
	consumers := make([]*sarama.Consumer, count)
	i := 0
	for partition, offset := range offsets {
		if consumer, err := c.getConsumer(c.client, c.Topic, int32(partition), c.ClientName, offset); err != nil {
			panic(err)
		} else {
			consumers[i] = consumer
		}
		i++
	}
	return consumers
}

func (c *Consumer) consume(consumers []*sarama.Consumer) {
	c.stopper = make(chan struct{})
	for _, consumer := range consumers {
		defer consumer.Close()
	}

	for _, consumer := range consumers {
		go func() {
			select {
			case value := <-consumer.Events():
				evt := &ConsumerEvent{
					Error:     value.Err,
					Message:   value.Value,
					Offset:    value.Offset,
					Partition: value.Partition,
				}
				c.events <- evt
			case <-c.stopper:
				return
			}
		}()
	}
}

func (c *Consumer) getConsumer(client *sarama.Client, topic string, partition int32, consumerName string, startAtOffset int64) (*sarama.Consumer, error) {
	config := sarama.NewConsumerConfig()
	config.EventBufferSize = 1
	config.OffsetMethod = sarama.OffsetMethodManual
	config.OffsetValue = startAtOffset
	consumer, err := sarama.NewConsumer(client, topic, partition, consumerName, config)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}
