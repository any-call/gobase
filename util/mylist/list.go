package mylist

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

type List[E comparable] struct {
	sync.RWMutex
	list []*node[E]
}

func New[E comparable]() *List[E] {
	return new(List[E]).init()
}

// 返回一个空数组
func (l *List[E]) init() *List[E] {
	l.list = make([]*node[E], 0, 500)
	return l
}

func (l *List[E]) Append(item E) {
	l.Lock()
	defer l.Unlock()

	l.list = append(l.list, newNode[E](item))
}

func (l *List[E]) PreAppend(item E) {
	l.Lock()
	defer l.Unlock()

	l.list = append([]*node[E]{newNode(item)}, l.list...)
}

func (l *List[E]) Insert(index int, item E) error {
	l.Lock()
	defer l.Unlock()

	if b := l.isValidIndex(index); !b {
		return errors.New("outside of range")
	}

	l.list = append(l.list[:index], append([]*node[E]{newNode[E](item)}, l.list[index:]...)...)
	return nil
}

func (l *List[E]) RemoveAt(index int) error {
	l.Lock()
	defer l.Unlock()

	if b := l.isValidIndex(index); !b {
		return errors.New("outside of range")
	}

	l.list = append(l.list[:index], append([]*node[E]{}, l.list[index+1:]...)...)
	return nil
}

func (l *List[E]) RemoveFirst() error {
	l.Lock()
	defer l.Unlock()

	if len(l.list) == 0 {
		return errors.New("empty range")
	}

	l.list = l.list[1:]
	return nil
}

func (l *List[E]) RemoveLast() error {
	l.Lock()
	defer l.Unlock()

	if len(l.list) == 0 {
		return errors.New("empty range")
	}

	l.list = l.list[:len(l.list)-1]
	return nil
}

func (l *List[E]) Move(from, to int) error {
	l.Lock()
	defer l.Unlock()

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

func (l *List[E]) At(index int) (v E, err error) {
	l.RLock()
	defer l.RUnlock()

	if b := l.isValidIndex(index); !b {
		err = errors.New("outside range")
	} else {
		v = l.list[index].value
	}

	return
}

func (l *List[E]) Len() int {
	l.RLock()
	defer l.RUnlock()

	return len(l.list)
}

func (l *List[E]) String() string {
	l.RLock()
	defer l.RUnlock()
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
