package kafka

import (
	"github.com/gotgo/fw/logging"
	//	"github.com/gotgo/fw/me"
	"github.com/kraveio/sarama"
)

type Client struct {
	Hosts     []string
	Log       logging.Logger `inject:""`
	asyncSend chan *sarama.Message
}

//func Connect(hosts []string) (*Client, error) {
//	client := &Client{
//		Hosts: hosts,
//	}
//
//	if sclient, err := client.doConnect(); err != nil {
//		return nil, err
//	} else {
//		client.client = sclient
//	}
//	return client, nil
//}

//func (c *Client) runProducer() {
//	producer, err := c.asyncProducer()
//	if err != nil {
//		panic(err)
//	}
//	defer producer.Close()
//	for {
//		select {
//		case message := <-c.asyncSend:
//			producer.Input() <- message
//		case err := <-producer.Errors():
//			//what do we do?
//			panic(err.Err)
//		}
//	}
//}
//
//func (c *Client) doConnect() (*sarama.Client, error) {
//	hosts := c.Hosts
//	if len(hosts) == 0 {
//		return nil, errors.New("no kafka hosts configured")
//	}
//	clientName, err := os.Hostname()
//	if err != nil {
//		clientName = "hostnameFail"
//	}
//	client, err := sarama.NewClient(hosts, sarama.NewClientConfig())
//	if err != nil {
//		return nil, errors.New("failed to connect to kafka")
//	}
//	return client, nil
//}
//
//func (c *Client) Close() {
//	client := c.client
//	if client != nil {
//		client.Close()
//		client = nil
//	}
//}

//func (c *Client) sendSync(bts []byte, topic, key string) error {
//	config := sarama.NewProducerConfig()
//	config.AckSuccesses = true
//	config.Partitioner = func() sarama.Partitioner { return sarama.NewHashPartitioner() }
//	producer, err := sarama.NewSimpleProducer(c.client, config)
//	if err != nil {
//		return me.Err(err, "failed to create kafka simple producer")
//	}
//	defer producer.Close()
//	if err = producer.SendMessage(topic, sarama.StringEncoder(key), sarama.ByteEncoder(bts)); err != nil {
//		return me.Err(err, "kafka send message fail")
//	}
//	return nil
//}

func (c *Client) send(bts []byte, topic, key string) error {
	config := sarama.NewConfig()
	producer, err := sarama.NewSyncProducer(c.Hosts, config)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
		}
	}()

	msg := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(bts)}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}

//func (c *Client) asyncProducer() (*sarama.Producer, error) {
//	client := c.client
//	if client == nil {
//		panic("must establish a connection before using a producer")
//	}
//
//	cfg := sarama.NewProducerConfig()
//	cfg.Partitioner = func() sarama.Partitioner { return sarama.NewHashPartitioner() }
//	cfg.FlushMsgCount = 1
//	cfg.FlushFrequency = time.Millisecond * 100
//
//	if producer, err := sarama.NewProducer(client, cfg); err != nil {
//		return nil, err
//	} else {
//		return producer, nil
//	}
//}

func (c *Client) SendBytesAsync(bts []byte, topic, key string) {
	//	encodedKey := sarama.StringEncoder(key)
	//	c.asyncSend <- &sarama.MessageToSend{Topic: topic, Key: encodedKey, Value: sarama.ByteEncoder(bts)}
	c.send(bts, topic, key)
}

func (c *Client) SendBytes(bts []byte, topic, key string) error {
	return c.send(bts, topic, key)
}

func (c *Client) NewConsumer(name, topic string, startingOffsets map[int32]int64) *Consumer {
	return NewConsumer(c.Hosts, name, topic, startingOffsets)
}
