package main

type IdMap_t[T any] struct {
	sort []int
	byId map[int]T
}

func IdMap[T any]() IdMap_t[T] {
	return IdMap_t[T]{
		sort: []int{},
		byId: make(map[int]T),
	}
}

func (in IdMap_t[T])Add(id int, item T) IdMap_t[T] {
	if _, ok := in.byId[id]; !ok { in.sort = append(in.sort, id) }
	in.byId[id] = item
	return in
}
