package mylist

type node[E any] struct {
	value E
}

func newNode[E any](v E) *node[E] {
	return new(node[E]).init(v)
}

func (e *node[E]) init(v E) *node[E] {
	e.value = v
	return e
}

func (e *node[E]) Value() E {
	return e.value
}
