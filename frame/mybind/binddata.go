package mybind

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

var shareBaseBindObj BindData = newBaseBind()

type baseBind struct {
	listeners    sync.Map //map[listener]map[data]true
	datas        sync.Map //map[data]map[listener]true
	dataVal      sync.Map // map[data]map[reflect.value]current value
	timeDuration time.Duration
}

func newBaseBind() BindData {
	ret := &baseBind{}
	go ret.trigger()
	return ret
}

func (b *baseBind) AddListener(listerer Listener, data any) error {
	//检测是不是有效的数据
	if listerer == nil {
		return fmt.Errorf("listner is nil")
	}

	if data == nil {
		return fmt.Errorf("val is nil")
	}

	if t := reflect.TypeOf(data); t.Kind() != reflect.Pointer {
		return fmt.Errorf("%v must point type", data)
	}

	pointer := reflect.ValueOf(data)
	newValue := pointer.Elem()

	if newValue.CanSet() == false {
		return fmt.Errorf("%v can't be changed", data)
	}

	//记录初始化
	if _, ok := b.dataVal.Load(data); !ok {
		valMap := &sync.Map{}
		valMap.Store(newValue, refValue(newValue))
		b.dataVal.Store(data, valMap)
	}

	//增加 监听 ==》数据
	if topVal, ok := b.listeners.Load(listerer); ok {
		if subMap, ok := topVal.(*sync.Map); ok {
			subMap.Store(data, true)
		} else {
			return fmt.Errorf("unexcepted type:%v", topVal)
		}
	} else {
		topVal := &sync.Map{}
		topVal.Store(data, true)
		b.listeners.Store(listerer, topVal)
	}

	//增加 数据 ==》listener
	if topVal, ok := b.datas.Load(data); ok {
		if subMap, ok := topVal.(*sync.Map); ok {
			subMap.Store(listerer, true)
		} else {
			return fmt.Errorf("unexcepted type:%v", topVal)
		}
	} else {
		topVal := &sync.Map{}
		topVal.Store(listerer, true)
		b.datas.Store(data, topVal)
	}

	return nil
}

func (b *baseBind) RemoteListener(listener Listener) error {
	var err error
	if topVal1, ok := b.listeners.LoadAndDelete(listener); ok {
		if dataMap, ok := topVal1.(*sync.Map); ok {
			dataMap.Range(func(data, _ any) bool {
				if data != nil {
					//清掉监听的相关数据引用
					if topVal2, ok := b.datas.Load(data); ok {
						if listenerMap, ok := topVal2.(*sync.Map); ok {
							listenerMap.Delete(listener)
							length := 0
							listenerMap.Range(func(_, _ any) bool {
								length++
								return true
							})

							if length == 0 {
								b.datas.Delete(data)
								b.dataVal.Delete(data)
							}
						} else {
							err = fmt.Errorf("unexcepted type:%v", topVal2)
							return false
						}
					}
				}
				return true
			})
		} else {
			return fmt.Errorf("unexcepted type:%v", topVal1)
		}
	}

	return err
}

func (b *baseBind) trigger() {
	for {
		b.dataVal.Range(func(data, value any) bool {
			if valueMap, ok := value.(*sync.Map); ok {
				valueMap.Range(func(key, oldVal any) bool {
					if reValue, ok := key.(reflect.Value); ok {
						newValue := refValue(reValue)
						bFlag := reflect.DeepEqual(oldVal, newValue)
						if !bFlag { //值已改变
							//fmt.Println("newValue :", newValue, ";oldValue :", oldVal)
							//首先存入新值
							valueMap.Store(key, newValue)

							//notification to all
							if valMap, ok := b.datas.Load(data); ok {
								if listenerMap, ok := valMap.(*sync.Map); ok {
									var listenerList []Listener = []Listener{}
									listenerMap.Range(func(key, _ any) bool {
										if v, ok := key.(Listener); ok {
											listenerList = append(listenerList, v)
										}
										return true
									})

									if len(listenerList) > 0 {
										go func(listers []Listener, data any) {
											for i, _ := range listers {
												go listers[i].DataChanged(data)
											}
										}(listenerList, data)
									}
								}
							}
						}
					}
					return true
				})
			}
			return true
		})

		if b.timeDuration > 0 {
			time.Sleep(b.timeDuration)
		}
	}
}

func refValue(v reflect.Value) any {
	return fmt.Sprintf("%v", v.Interface())
}
