package set

type Set[T comparable] map[T]struct{}

func (set Set[T]) Add(s T) {
	set[s] = struct{}{}
}

func (set Set[T]) Remove(s T) {
	delete(set, s)
}

func (set Set[T]) Has(s T) bool {
	_, ok := set[s]
	return ok
}

func (set Set[T]) Len() int {
	return len(set)
}

func (set Set[T]) List() []T {
	var slice []T
	for s := range set {
		slice = append(slice, s)
	}
	return slice
}

func (set Set[T]) ToMap() map[T]struct{} {
	return map[T]struct{}(set)
}

func (set Set[T]) Union(other Set[T]) Set[T] {
	union := Set[T]{}
	for s := range set {
		union.Add(s)
	}
	for s := range other {
		union.Add(s)
	}
	return union
}

func (set Set[T]) Iterate(f func(T) bool) {
	for s := range set {
		if !f(s) {
			break
		}
	}
}

func (set Set[T]) IterateAll(f func(T)) {
	for s := range set {
		f(s)
	}
}

func New[T comparable](strs ...T) Set[T] {
	set := Set[T]{}
	for _, s := range strs {
		set.Add(s)
	}
	return set
}
