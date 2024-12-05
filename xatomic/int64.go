package xatomic

import "sync/atomic"

func DecreaseInt64WithoutBelowZero(value *int64) {
	for {
		current := atomic.LoadInt64(value)
		if current == 0 {
			return
		}

		if atomic.CompareAndSwapInt64(value, current, current-1) {
			return
		}
	}
}
