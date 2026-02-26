package mylog

// Trace 是一次调用链的上下文（不可变链式结构）
type Trace struct {
	parent *Trace
	key    string
	value  any
}

// New 创建根 Trace
func NewTrace() *Trace {
	return nil
}

// With 增加一个字段（零 map 拷贝）
func (t *Trace) With(key string, value any) *Trace {
	return &Trace{
		parent: t,
		key:    key,
		value:  value,
	}
}

// WithKV 批量增加字段
func (t *Trace) WithKV(kv ...any) *Trace {
	n := t
	for i := 0; i < len(kv)-1; i += 2 {
		k, _ := kv[i].(string)
		n = n.With(k, kv[i+1])
	}
	return n
}

func (t *Trace) toKV() []any {
	if t == nil {
		return nil
	}

	// 收集链
	stack := make([]*Trace, 0, 8) // 小容量起步
	for cur := t; cur != nil; cur = cur.parent {
		stack = append(stack, cur)
	}

	kv := make([]any, 0, len(stack)*2)

	// 逆序写入
	for i := len(stack) - 1; i >= 0; i-- {
		if stack[i].key != "" {
			kv = append(kv, stack[i].key, stack[i].value)
		}
	}

	return kv
}

func (t *Trace) Info(event string, kv ...any) {
	InfoKV(event, append(t.toKV(), kv...)...)
}

func (t *Trace) Debug(event string, kv ...any) {
	DebugKV(event, append(t.toKV(), kv...)...)
}

func (t *Trace) Warn(event string, kv ...any) {
	WarnKV(event, append(t.toKV(), kv...)...)
}
