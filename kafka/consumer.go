package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/gotgo/fw/logging"
)

type Consumer struct {
	Name            string
	StartingOffsets map[int]int64
	Log             logging.Logger `inject:""`
	Topic           string
	Factory         func(topic string, partition int32, consumerName string, startAtOffset int64) (ConsumerChannel, error)
	stopper         chan struct{}
	events          chan *ConsumerEvent
}

func NewConsumer(name, topic string, startingOffsets map[int]int64, factory func(topic string, partition int32, consumerName string, startAtOffset int64) (ConsumerChannel, error)) (*Consumer, error) {
	consumer := new(Consumer)
	consumer.Name = name
	consumer.StartingOffsets = startingOffsets
	consumer.Topic = topic
	consumer.Factory = factory
	if err := consumer.consumePartitions(); err != nil {
		return nil, err
	} else {
		return consumer, nil
	}
}

func (c *Consumer) Events() <-chan *ConsumerEvent {
	return c.events
}

func (c *Consumer) Close() {
	close(c.stopper)
}

func (c *Consumer) consumePartitions() error {
	c.stopper = make(chan struct{})
	offsets := c.StartingOffsets
	count := len(offsets)
	consumers := make([]ConsumerChannel, count, count)
	i := 0
	for partition, offset := range offsets {
		if consumer, err := c.Factory(c.Topic, int32(partition), c.Name, int64(offset)); err != nil {
			return err
		} else {
			consumers[i] = consumer
		}
		i++
	}
	go c.consume(consumers)
	return nil
}

func (c *Consumer) consume(consumers []ConsumerChannel) {
	for _, consumer := range consumers {
		defer consumer.Close()
	}

	for {
		for _, consumer := range consumers {
			select {
			case <-c.stopper:
				return
			case event := <-consumer.Events():
				if event.Err != nil {
					c.Log.Error("error retrieving message from kafka", event.Err)
				} else {
					c.events <- event
				}
			default:
				continue
			}
		}
	}
}

type consumerWrapper struct {
	Consumer *sarama.Consumer
	Channel  chan *ConsumerEvent
	Stopper  chan struct{}
}

func (cw *consumerWrapper) Open() {
	cw.Channel = make(chan *ConsumerEvent)
	cw.Stopper = make(chan struct{})
	for {
		select {
		case value := <-cw.Consumer.Events():
			//return wrapped event
			evt := &ConsumerEvent{
				Err:       value.Err,
				Message:   value.Value,
				Offset:    value.Offset,
				Partition: value.Partition,
			}
			cw.Channel <- evt
		case <-cw.Stopper:
			return
		}
	}
}

func (cw *consumerWrapper) Events() <-chan *ConsumerEvent {
	return cw.Channel
}

func (cw *consumerWrapper) Close() error {
	close(cw.Stopper)
	return cw.Consumer.Close()
}
