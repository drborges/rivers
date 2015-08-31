package stream

import (
	"sort"
)

func (by SortByFn) Sort(items []T) {
	sort.Sort(&sorter{
		items: items,
		by:    by,
	})
}

type sorter struct {
	items []T
	by    SortByFn
}

func (s *sorter) Len() int {
	return len(s.items)
}

func (s *sorter) Swap(i, j int) {
	s.items[i], s.items[j] = s.items[j], s.items[i]
}

func (s *sorter) Less(i, j int) bool {
	return s.by(s.items[i], s.items[j])
}
