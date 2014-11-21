package kafka

import (
	"errors"
	"os"

	"github.com/Shopify/sarama"
	"github.com/gotgo/fw/logging"
)

type Client struct {
	Hosts  []string
	Log    logging.Logger `inject:""`
	client *sarama.Client
}

func Connect(hosts []string) (*Client, error) {
	client := &Client{
		Hosts: hosts,
	}

	if sclient, err := client.doConnect(); err != nil {
		return nil, err
	} else {
		client.client = sclient
	}
	return client, nil
}

func (c *Client) doConnect() (*sarama.Client, error) {
	hosts := c.Hosts
	if len(hosts) == 0 {
		return nil, errors.New("no kafka hosts configured")
	}
	clientName, err := os.Hostname()
	if err != nil {
		clientName = "hostnameFail"
	}
	client, err := sarama.NewClient(clientName, hosts, sarama.NewClientConfig())
	if err != nil {
		return nil, errors.New("failed to connect to kafka")
	}
	return client, nil
}

func (c *Client) Close() {
	client := c.client
	if client != nil {
		client.Close()
		client = nil
	}
}

func (c *Client) producer() (*sarama.Producer, error) {
	client := c.client
	if client == nil {
		panic("must establish a connection before using a producer")
	}

	cfg := sarama.NewProducerConfig()
	cfg.Partitioner = sarama.NewHashPartitioner()
	cfg.MaxBufferedBytes = 2
	cfg.MaxBufferTime = 2

	if producer, err := sarama.NewProducer(client, cfg); err != nil {
		return nil, err
	} else {
		return producer, nil
	}
}

func (c *Client) SendBytes(message []byte, topic, key string) error {
	producer, err := c.producer()
	if err != nil {
		return err
	}
	defer producer.Close()

	encodedKey := sarama.StringEncoder(key)
	err = producer.SendMessage(topic, encodedKey, sarama.ByteEncoder(message))
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) SendString(message, topic, key string) error {
	producer, err := c.producer()
	if err != nil {
		return err
	}
	defer producer.Close()

	encodedKey := sarama.StringEncoder(key)
	payload := sarama.StringEncoder(message)
	err = producer.SendMessage(topic, encodedKey, payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) NewConsumer(name, topic string, startingOffsets map[PartitionIndex]Offset) (*Consumer, error) {
	return NewConsumer(name, topic, startingOffsets, c.consumerFactory)
}

func (c *Client) consumerFactory(topic string, partition int32, consumerName string, startAtOffset int64) (ConsumerChannel, error) {
	config := sarama.NewConsumerConfig()
	config.EventBufferSize = 2
	config.OffsetMethod = sarama.OffsetMethodManual
	config.OffsetValue = startAtOffset
	if consumer, err := sarama.NewConsumer(c.client, topic, partition, consumerName, config); err != nil {
		return nil, err
	} else {
		wrapper := &consumerWrapper{
			Consumer: consumer,
		}
		go wrapper.Open()
		return wrapper, nil
	}
}
