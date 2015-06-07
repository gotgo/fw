package kafka

import (
	"sync"

	"github.com/gotgo/fw/logging"
	"github.com/kraveio/sarama"
)

type Consumer struct {
	ClientName      string
	StartingOffsets map[int32]int64
	Log             logging.Logger `inject:""`
	Topic           string
	stopper         chan struct{}
	events          chan *ConsumerEvent
	hosts           []string
}

func NewConsumer(hosts []string, clientName, topic string, startingOffsets map[int32]int64) *Consumer {
	consumer := &Consumer{
		ClientName:      clientName,
		StartingOffsets: startingOffsets,
		Topic:           topic,
		events:          make(chan *ConsumerEvent),
		hosts:           hosts,
	}
	return consumer
}

// Run - Nonblocking
func (c *Consumer) Run() {
	consumers := c.getConsumerChannels()
	go c.consume(consumers)
}

func (c *Consumer) Events() <-chan *ConsumerEvent {
	return c.events
}

func (c *Consumer) Shutdown() {
	if c.stopper != nil {
		close(c.stopper)
	}
}

func (c *Consumer) getConsumerChannels() []<-chan *sarama.ConsumerMessage {
	offsets := c.StartingOffsets
	consumers := make([]<-chan *sarama.ConsumerMessage, len(offsets))
	i := 0
	for partition, offset := range offsets {
		consumer, err := c.getConsumer(c.hosts, c.Topic, int32(partition), c.ClientName, offset)

		if err != nil {
			panic(err)
		}

		consumers[i] = consumer.Messages()
		i++
	}
	return consumers
}

func (c *Consumer) consume(consumers []<-chan *sarama.ConsumerMessage) {
	stopper := make(chan struct{})
	single := merge(stopper, consumers...)
	c.stopper = stopper

	for value := range single {
		evt := &ConsumerEvent{
			Error:     nil,
			Message:   value.Value,
			Offset:    value.Offset,
			Partition: value.Partition,
		}
		c.events <- evt
	}
}

func merge(done <-chan struct{}, cs ...<-chan *sarama.ConsumerMessage) <-chan *sarama.ConsumerMessage {
	var wg sync.WaitGroup
	out := make(chan *sarama.ConsumerMessage)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *sarama.ConsumerMessage) {
		defer wg.Done()
		for e := range c {
			select {
			case out <- e:
			case <-done:
				return
			}
		}
	}
	wg.Add(len(cs))
	for _, c := range cs {
		local := c
		go output(local)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (c *Consumer) getConsumer(hosts []string, topic string, partition int32, consumerName string, startAtOffset int64) (sarama.PartitionConsumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(hosts, config)
	if err != nil {
		return nil, err
	}
	return consumer.ConsumePartition(topic, partition, startAtOffset)
}
