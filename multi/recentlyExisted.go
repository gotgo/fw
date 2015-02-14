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
			u := r.current - 1
			if u < 0 {
				u = 0
			}
			r.list[u] = t
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
