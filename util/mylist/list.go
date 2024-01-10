package mylist

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type List[E any] struct {
	lock        sync.RWMutex
	list        []*node[E]
	capacityNum int32
}

func NewList[E any]() *List[E] {
	return new(List[E]).init(500)
}

func NewListEx[E any](capacityNum int32) *List[E] {
	return new(List[E]).init(capacityNum)
}

// 返回一个空数组
func (l *List[E]) init(num int32) *List[E] {
	l.list = make([]*node[E], 0, num)
	l.capacityNum = num
	return l
}

func (l *List[E]) Append(item E) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.list = append(l.list, newNode[E](item))
}

func (l *List[E]) AppendByArray(items []E) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if items != nil {
		nodes := make([]*node[E], len(items))
		for i, _ := range items {
			nodes[i] = newNode(items[i])
		}

		l.list = append(l.list, nodes...)
	}
}

func (l *List[E]) ResetByArray(items []E) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.list = make([]*node[E], 0, 500)
	if items != nil {
		for i, _ := range items {
			l.list = append(l.list, newNode(items[i]))
		}
	}

	return
}

func (l *List[E]) PreAppend(item E) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.list = append([]*node[E]{newNode(item)}, l.list...)
}

func (l *List[E]) Insert(index int, item E) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if b := l.isValidIndex(index); !b {
		return errors.New("outside of range")
	}

	l.list = append(l.list[:index], append([]*node[E]{newNode[E](item)}, l.list[index:]...)...)
	return nil
}

func (l *List[E]) RemoveAt(index int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if b := l.isValidIndex(index); !b {
		return errors.New("outside of range")
	}

	l.list = append(l.list[:index], l.list[index+1:]...)
	return nil
}

func (l *List[E]) RemoveFirst() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.list) == 0 {
		return errors.New("empty range")
	}

	l.list = l.list[1:]
	return nil
}

func (l *List[E]) RemoveLast() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.list) == 0 {
		return errors.New("empty range")
	}

	l.list = l.list[:len(l.list)-1]
	return nil
}

func (l *List[E]) Clear() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.init(l.capacityNum)

	return
}

func (l *List[E]) Move(from, to int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.isValidIndex(from) == false {
		return errors.New("invalid index")
	}
	if l.isValidIndex(to) == false {
		return errors.New("invalid index")
	}

	if from == to {
		return nil
	}

	if from < to {
		backList := l.list[to+1:]
		if len(backList) > 0 {
			backList = append([]*node[E]{newNode[E](l.list[to+1].value)}, backList[1:]...)
		}
		l.list = append(append(l.list[:from], append(l.list[from+1:to+1], l.list[from])...), backList...)
	} else {
		toItem := newNode[E](l.list[to].value)
		toFromList := append([]*node[E]{toItem}, l.list[to+1:from]...)
		toFromList = append(toFromList, l.list[from+1:]...)
		l.list = append(append(l.list[:to], l.list[from]), toFromList...)
	}
	return nil
}

func (l *List[E]) SwapItemsAt(idx1, idx2 int) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.isValidIndex(idx1) == false {
		return errors.New("invalid index")
	}
	if l.isValidIndex(idx2) == false {
		return errors.New("invalid index")
	}

	if idx1 == idx2 {
		return nil
	}

	tmpNode := l.list[idx1]
	l.list[idx1] = l.list[idx2]
	l.list[idx2] = tmpNode
	return nil
}

func (l *List[E]) First() (v E, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if len(l.list) == 0 {
		err = errors.New("empty list")
		return
	}

	v = l.list[0].value
	return
}

func (l *List[E]) Last() (v E, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if len(l.list) == 0 {
		err = errors.New("empty list")
		return
	}

	v = l.list[len(l.list)-1].value
	return
}

func (l *List[E]) At(index int) (v E, err error) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if b := l.isValidIndex(index); !b {
		err = errors.New("outside range")
	} else {
		v = l.list[index].value
	}

	return
}

func (l *List[E]) Len() int {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return len(l.list)
}

func (l *List[E]) TakeFirst() (v E, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.list) == 0 {
		err = errors.New("empty list")
		return
	}

	v = l.list[0].value
	l.list = l.list[1:]
	return
}

func (l *List[E]) TakeLast() (v E, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if len(l.list) == 0 {
		err = errors.New("empty list")
		return
	}

	v = l.list[len(l.list)-1].value
	l.list = l.list[:len(l.list)-1]
	return
}

func (l *List[E]) TakeAt(idx int) (v E, err error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	if b := l.isValidIndex(idx); !b {
		err = errors.New("invalid index")
		return
	}

	v = l.list[idx].value
	l.list = append(l.list[:idx], l.list[idx+1:]...)
	return
}

func (l *List[E]) TakeAll() []E {
	l.lock.Lock()
	defer l.lock.Unlock()

	ret := make([]E, len(l.list))

	for i, _ := range l.list {
		ret[i] = l.list[i].value
	}

	l.init(l.capacityNum)
	return ret
}

func (l *List[E]) TakeHeadN(count int) []E {
	l.lock.Lock()
	defer l.lock.Unlock()

	if count <= 0 {
		return []E{}
	}

	if count > len(l.list) {
		count = len(l.list)
	}

	ret := make([]E, count)

	for i := 0; i < count; i++ {
		ret[i] = l.list[i].value
	}

	l.list = l.list[count:]

	return ret
}

func (l *List[E]) Range(f func(index int, v E)) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for i, v := range l.list {
		f(i, v.value)
	}
}

func (l *List[E]) Where(f func(index int, v E) bool) *List[E] {
	l.lock.RLock()
	defer l.lock.RUnlock()

	ret := NewList[E]()
	for i, v := range l.list {
		if f(i, v.value) {
			ret.Append(v.value)
		}
	}

	return ret
}

func (l *List[E]) ToArray() []E {
	l.lock.RLock()
	defer l.lock.RUnlock()
	ret := make([]E, len(l.list))
	for i, _ := range l.list {
		ret[i] = l.list[i].Value()
	}
	return ret
}

func (l *List[E]) String() string {
	l.lock.RLock()
	defer l.lock.RUnlock()
	tmpList := make([]string, len(l.list))
	for i, _ := range l.list {
		tmpList[i] = fmt.Sprintf("%v", l.list[i].Value())
	}
	return strings.Join(tmpList, ";")
}

// inner function
func (l *List[E]) isValidIndex(index int) bool {
	if index >= 0 && index < len(l.list) {
		return true
	}

	return false
}
