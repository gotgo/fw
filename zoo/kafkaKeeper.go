package zoo

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gotgo/fw/me"
	"github.com/samuel/go-zookeeper/zk"
)

const connectTimeout = 3 * time.Second

func NewKafkaKeeper(hosts []string, c *TopicConsumer, s *KafkaState) *KafkaKeeper {
	return &KafkaKeeper{
		acl:      zk.WorldACL(zk.PermAll),
		hosts:    hosts,
		Consumer: c,
		State:    s,
	}
}

// KafkaKeeper - Used to manage consumer per topic, per application consuming
// > /{TopicConsumer.Root}/{topic}/apps/{app}/partitions
//														/{partition}/consumed/{offset}
type KafkaKeeper struct {
	acl      []zk.ACL
	hosts    []string
	Consumer *TopicConsumer
	State    *KafkaState
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

		current := path.Join(current, part)
		fmt.Printf("cheking for %s", current)
		exists, _, err := c.Exists(current)
		if err != nil {
			return me.Err(err, "error checking if path exists: "+current)
		}
		if !exists {
			const flags = 0
			fmt.Printf("creating %s", current)
			fmt.Println()
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

	//  /kafka-topics/{topic}/apps/{app}/partitions/"
	if err := z.ensureExists(c, tc.PartitionsPath(), ""); err != nil {
		return err
	}

	return nil
}

func (z *KafkaKeeper) GetOffsets() (map[int]int64, error) {
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

	offsets := make(map[int]int64)
	for _, p := range partitions {
		//{partition}/consumed/{offset}
		err := z.ensureExists(conn, path.Join(z.Consumer.PartitionsPath(), p, "consumed"), "0")
		if err != nil {
			return nil, err
		}
		offset, _, err := z.get(conn, path.Join(z.Consumer.PartitionsPath(), p, "consumed"))
		if err != nil {
			return nil, err
		}
		k, _ := strconv.Atoi(p)
		v, _ := strconv.ParseInt(offset, 10, 0)
		offsets[k] = v
	}
	return offsets, nil
}

func (z *KafkaKeeper) SetOffset(partition int32, offset int64) error {
	conn := z.connect()
	defer conn.Close()

	err := z.ensureSetup(conn)
	if err != nil {
		return me.Err(err, fmt.Sprintf("failed to set offset on partition: %d", partition))
	}
	path := path.Join(z.Consumer.PartitionsPath(), fmt.Sprintf("%d", partition), "consumed")

	//force version for now
	_, version, err := z.get(conn, path)
	if err != nil {
		return err
	}

	err = z.set(conn, path, fmt.Sprintf("%d", offset), version)
	if err != nil {
		return me.Err(err, "failed to set value for path: "+path)
	}
	return err
}
