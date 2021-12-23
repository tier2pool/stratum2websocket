package pool

import "sync"

const (
	DefaultSize = 2048
)

var pool = sync.Pool{
	New: func() interface{} {
		return make([]byte, DefaultSize)
	},
}

func Get() []byte {
	return pool.Get().([]byte)
}

func Put(b []byte) {
	if len(b) == DefaultSize {
		pool.Put(b)
	}
}
