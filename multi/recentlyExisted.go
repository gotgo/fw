package multi

func NewRecentlyExisted(size int) *RecentlyExisted {
	return &RecentlyExisted{
		list: make([]string, size),
	}
}

type RecentlyExisted struct {
	list    []string
	current int
}

func (r *RecentlyExisted) CheckAndAdd(t string) bool {
	for _, c := range r.list {
		if c == t {
			return true
		}
	}

	r.list[r.current] = t
	r.current++
	if r.current == len(r.list) {
		r.current = 0
	}
	return false
}
