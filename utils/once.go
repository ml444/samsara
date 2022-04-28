package utils

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once)Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.m.Lock()
		defer o.m.Unlock()
		if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
			f()
		}
	}
}
	
	
