package kafka

import (
	"sync"

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

func (c *Consumer) getConsumerChannels() []<-chan *sarama.ConsumerEvent {
	offsets := c.StartingOffsets
	consumers := make([]<-chan *sarama.ConsumerEvent, len(offsets))
	i := 0
	for partition, offset := range offsets {
		consumer, err := c.getConsumer(c.client, c.Topic, int32(partition), c.ClientName, offset)

		if err != nil {
			panic(err)
		}

		consumers[i] = consumer.Events()
		i++
	}
	return consumers
}

func (c *Consumer) consume(consumers []<-chan *sarama.ConsumerEvent) {
	stopper := make(chan struct{})
	single := merge(stopper, consumers...)
	c.stopper = stopper

	for value := range single {
		evt := &ConsumerEvent{
			Error:     value.Err,
			Message:   value.Value,
			Offset:    value.Offset,
			Partition: value.Partition,
		}
		c.events <- evt
	}
}

func merge(done <-chan struct{}, cs ...<-chan *sarama.ConsumerEvent) <-chan *sarama.ConsumerEvent {
	var wg sync.WaitGroup
	out := make(chan *sarama.ConsumerEvent)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan *sarama.ConsumerEvent) {
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
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func (c *Consumer) getConsumer(client *sarama.Client, topic string, partition int32, consumerName string, startAtOffset int64) (*sarama.Consumer, error) {
	config := sarama.NewConsumerConfig()
	config.EventBufferSize = 1
	config.OffsetMethod = sarama.OffsetMethodManual
	config.OffsetValue = startAtOffset
	return sarama.NewConsumer(client, topic, partition, consumerName, config)
}
