package mutli

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
			u := current - 1
			if u < 0 {
				u = 0
			}
			list[u] = t
			return true
		}
	}

	list[current] = t
	current++
	if current == len(list) {
		current = 0
	}
	return false
}
