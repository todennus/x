package xbytes

import (
	"math"
	"sync"
	"sync/atomic"

	"github.com/todennus/x/xatomic"
)

type bufferPoolAction int

const (
	bufferPoolActionPut bufferPoolAction = iota
	bufferPoolActionGet
)

type sizedPool struct {
	sync.Pool
	numAvailableObjects int64
}

func (p *sizedPool) Size() int64 {
	return atomic.LoadInt64(&p.numAvailableObjects)
}

var bufferPoolMutex = sync.RWMutex{}
var bufferPool = map[int]*sizedPool{}

func getBufferPoolByFactor(factor int, action bufferPoolAction) *sizedPool {
	if bufferPool[factor] == nil {
		bufferPoolMutex.Lock()

		if bufferPool[factor] == nil {
			capability := int(math.Pow(2, float64(factor)))

			bufferPool[factor] = &sizedPool{
				Pool: sync.Pool{
					New: func() any {
						return make([]byte, 0, capability)
					},
				},
				numAvailableObjects: 0,
			}
		}

		bufferPoolMutex.Unlock()
	}

	if action == bufferPoolActionPut {
		atomic.AddInt64(&bufferPool[factor].numAvailableObjects, 1)
	} else if action == bufferPoolActionGet {
		xatomic.DecreaseInt64WithoutBelowZero(&bufferPool[factor].numAvailableObjects)
	}

	return bufferPool[factor]
}

func getBufferPool(cap int64, action bufferPoolAction) *sizedPool {
	return getBufferPoolByFactor(int(math.Log2(float64(cap))), action)
}

func GetMostUsedBuffer() []byte {
	mostUsedFactor := -1
	mostPoolObjects := 0

	bufferPoolMutex.RLock()
	for k, v := range bufferPool {
		poolSize := int(v.Size())
		if mostPoolObjects < poolSize {
			mostUsedFactor = k
			mostPoolObjects = poolSize
		}
	}
	bufferPoolMutex.RUnlock()

	if mostUsedFactor == -1 {
		mostUsedFactor = 10 // default by 2^10 capability buffer
	}

	return getBufferPoolByFactor(mostUsedFactor, bufferPoolActionGet).Get().([]byte)
}

// GetBuffer returns a buffer with empty length and given cap from pool.
func GetBuffer(cap int64) []byte {
	return getBufferPool(cap, bufferPoolActionGet).Get().([]byte)
}

// PutBuffer resets the underlying of the buffer and adds it to the pool for
// future reuses.
func PutBuffer(buf []byte) {
	buf = buf[:0] // reset length of buffer
	getBufferPool(int64(cap(buf)), bufferPoolActionPut).Put(buf)
}
