package zoo

import (
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gotgo/fw/me"
	"github.com/samuel/go-zookeeper/zk"
)

const (
	connectTimeout        = 3 * time.Second
	defaultKafkaRoot      = "kafka"
	defaultKafkaTopicRoot = "kafka-topics"
)

func NewKafkaKeeper(hosts []string, c *TopicConsumer, s *KafkaState) *KafkaKeeper {
	if c.Root == "" {
		c.Root = defaultKafkaTopicRoot
	}
	if s.Root == "" {
		s.Root = defaultKafkaRoot
	}
	return &KafkaKeeper{
		acl:      zk.WorldACL(zk.PermAll),
		hosts:    hosts,
		Consumer: c,
		State:    s,
	}
}

// KafkaKeeper - Used to manage consumer per topic, per application consuming
// > /{TopicConsumer.Root}/{topic}/consumers/{consumer-app}/partitions
//														/{partition}/consumed/{offset}
type KafkaKeeper struct {
	acl      []zk.ACL
	hosts    []string
	Consumer *TopicConsumer
	State    *KafkaState
	mtx      sync.Mutex
}

func TestConnect(hosts []string) error {
	if conn, _, err := zk.Connect(hosts, connectTimeout); err != nil {
		return err
	} else {
		conn.Close()
		return nil
	}
}

func (k *KafkaKeeper) connect() *zk.Conn {
	k.mtx.Lock()
	defer k.mtx.Unlock()
	conn, _, err := zk.Connect(k.hosts, connectTimeout)
	if err != nil {
		panic(me.Err(err, "failed to connect to zookeeper nodes"))
	}
	return conn
}

func (z *KafkaKeeper) ensureExists(c *zk.Conn, p string, data string) error {
	parts := strings.Split(p, "/")
	current := "/"
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		current = path.Join(current, part)
		exists, _, err := c.Exists(current)
		if err != nil {
			return me.Err(err, "error checking if path exists: "+current)
		}
		if !exists {
			const flags = 0
			if _, err := c.Create(current, []byte(data), flags, z.acl); err != nil {
				return me.Err(err, "error creating path: "+current)
			}
		}
	}

	return nil
}

func (z *KafkaKeeper) get(c *zk.Conn, path string) (string, int32, error) {
	bts, stat, err := c.Get(path)
	if err != nil {
		return "", 0, me.Err(err, "failed to get path: "+path)
	}
	return fmt.Sprintf("%s", bts), stat.Version, nil
}

func (z *KafkaKeeper) set(c *zk.Conn, path, data string, version int32) error {
	if _, err := c.Set(path, []byte(data), version); err != nil {
		return me.Err(err, "failed to set path: "+path)
	}
	return nil
}

func (z *KafkaKeeper) getChildren(c *zk.Conn, path string) ([]string, error) {
	children, _, err := c.Children(path)
	if err != nil {
		return nil, me.Err(err, "failed to get children for path: "+path)
	}
	return children, nil
}

func (z *KafkaKeeper) ensureSetup(c *zk.Conn) error {
	tc := z.Consumer

	//  /kafka-topics
	if err := z.ensureExists(c, tc.Root, ""); err != nil {
		return err
	}

	//  /kafka-topics/{topic}/consumers/{consumer}/partitions/"
	if err := z.ensureExists(c, tc.PartitionsPath(), ""); err != nil {
		return err
	}

	return nil
}

// GetOffsets - returns an array sorted by Partition ascending
func (z *KafkaKeeper) GetOffsets() ([]*PartitionOffset, error) {
	conn := z.connect()
	defer conn.Close()

	err := z.ensureSetup(conn)
	if err != nil {
		return nil, me.Err(err, "failed to get offsets")
	}

	// get offical partition data from kafka settings
	partitions, err := z.getChildren(conn, z.State.PartitionsPath())
	if err != nil {
		return nil, me.Err(err, "failed to get kafka partition info")
	}

	po := make([]*PartitionOffset, len(partitions))
	for i, p := range partitions {
		//{partition}/consumed/{offset}
		err := z.ensureExists(conn, path.Join(z.Consumer.PartitionsPath(), p, "consumed"), "0")
		if err != nil {
			return nil, err
		}
		offset, _, err := z.get(conn, path.Join(z.Consumer.PartitionsPath(), p, "consumed"))
		if err != nil {
			return nil, err
		}
		p, _ := strconv.Atoi(p)
		o, _ := strconv.ParseInt(offset, 10, 0)
		po[i] = &PartitionOffset{Partition: int32(p), Offset: o}
	}
	sort.Sort(ByPartition{po})
	return po, nil
}

func (z *KafkaKeeper) SetOffset(partition int32, offset int64) error {
	conn := z.connect()
	defer conn.Close()

	err := z.ensureSetup(conn)
	if err != nil {
		return me.Err(err, fmt.Sprintf("failed to set offset on partition: %d", partition))
	}
	path := path.Join(z.Consumer.PartitionsPath(), fmt.Sprintf("%d", partition), "consumed")

	fmt.Printf("zookeeper path", path)
	//force version for now
	existingValue, version, err := z.get(conn, path)
	if err != nil {
		return err
	}

	//don't go backwards
	if eo, err := strconv.ParseInt(existingValue, 10, 64); err == nil && eo >= offset {
		return nil
	}

	err = z.set(conn, path, fmt.Sprintf("%d", offset), version)
	if err != nil {
		return me.Err(err, "failed to set value for path: "+path)
	}
	return err
}
