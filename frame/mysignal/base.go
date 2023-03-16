package mysignal

import (
	"sync"
)

func NewSignal[E comparable]() *Signal[E] {
	return &Signal[E]{}
}

// 定义 Signal结构
type Signal[E comparable] struct {
	sync.RWMutex
	cons  []con[E] //标识该信号联接的槽
	conID int      //记录每次联接的增量值
}

// 同步发射 Signal
func (sig *Signal[E]) Emit(param E) error {
	if sig == nil {
		return nil
	}

	sig.RLock()
	defer sig.RUnlock()

	if sig.cons == nil || len(sig.cons) == 0 {
		return nil
	}

	var err error
	wg := &sync.WaitGroup{}
	wg.Add(len(sig.cons))
	for i, _ := range sig.cons {
		go func(index int, info E) {
			defer wg.Done()
			if errTmp := sig.cons[index].slot(info); errTmp != nil {
				err = errTmp
			}
		}(i, param)
	}
	wg.Wait()
	return err
}

// 异步发射 Signal
func (sig *Signal[E]) AsyncEmit(param E) {
	go func(E) {
		sig.Emit(param)
	}(param)
}

// 增加 slot
func (sig *Signal[E]) AddSlot(s Slot[E]) int {
	if sig == nil {
		return 0
	}

	if s == nil {
		return 0
	}

	sig.RWMutex.Lock()
	defer sig.RWMutex.Unlock()

	sig.conID++
	newConntion := con[E]{
		slot: s,
		cid:  sig.conID,
	}

	sig.cons = append(sig.cons, newConntion)
	return sig.conID
}

// 删除  slot
func (sig *Signal[E]) DelSlot(cid int) {
	if sig == nil {
		return
	}

	if cid <= 0 || sig.cons == nil || len(sig.cons) == 0 {
		return
	}

	sig.RWMutex.Lock()
	defer sig.RWMutex.Unlock()

	for i, _ := range sig.cons {
		if sig.cons[i].cid == cid {
			sig.cons = append(sig.cons[:i], sig.cons[i+1:]...)
			return
		}
	}
}

type con[E comparable] struct {
	slot Slot[E]
	cid  int
}

// 定义 Slot
type Slot[E comparable] func(E) error

func Connect[E comparable](sig *Signal[E], s Slot[E]) int {
	if sig == nil {
		return 0
	}

	return sig.AddSlot(s)
}

func DisConnect[E comparable](sig *Signal[E], cid int) {
	if sig == nil {
		return
	}

	sig.DelSlot(cid)
}
