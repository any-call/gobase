package myenum

type enum[V any] struct {
	value V
	name  string
}

func NewENum[V any](name string, value V) ENum[V] {
	return &enum[V]{value: value, name: name}
}

func (self *enum[V]) Name() string {
	return self.name
}

func (self *enum[V]) Value() V {
	return self.value
}

func (self *enum[V]) SetValue(value V) {
	self.value = value
}
