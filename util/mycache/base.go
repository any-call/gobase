package mycache

import "time"

type Cache interface {
	Get(k string) (any, bool)             //取值
	Set(k string, v any, d time.Duration) //设置值
	HasKey(k string) bool
	UpdateExpired(k string, d time.Duration) error //更新值 时效
	SaveFile(fname string) error
	LoadFile(fname string) error
	Flush()
}

type Item struct {
	Object     any
	Expiration int64 //表示用效的时间Nano ,0表示没有到期时间
}

func (self Item) Expired() bool {
	if self.Expiration == 0 {
		return false
	}

	return time.Now().UnixNano() > self.Expiration
}
