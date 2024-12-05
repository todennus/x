package xbytes

import "sync"

var bytesPoolMutex = sync.Mutex{}
var bytesPool = map[int64]*sync.Pool{}

func getBytesPool(size int64) *sync.Pool {
	if bytesPool[size] == nil {
		bytesPoolMutex.Lock()

		if bytesPool[size] == nil {
			bytesPool[size] = &sync.Pool{
				New: func() any {
					return make([]byte, size)
				},
			}
		}

		bytesPoolMutex.Unlock()
	}

	return bytesPool[size]
}

// GetBytes returns a byte slice from pool.
func GetBytes(size int64) []byte {
	return getBytesPool(size).Get().([]byte)
}

// PutBytes adds the byte slice to the pool for future reuses.
func PutBytes(buf []byte) {
	getBytesPool(int64(len(buf))).Put(buf)
}
