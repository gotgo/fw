package zoo

type PartitionOffset struct {
	Partition int
	Offset    int64
}

type Offsets []*PartitionOffset

func (o Offsets) Len() int      { return len(o) }
func (o Offsets) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

type ByPartition struct{ Offsets }

func (p ByPartition) Less(i, j int) bool {
	return p.Offsets[i].Partition < p.Offsets[j].Partition
}
