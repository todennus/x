package xcrypto

import (
	"errors"
	"hash"
	"io"
)

var _ io.Reader = (*HashReader)(nil)

type HashReader struct {
	inner  io.Reader
	hasher hash.Hash
}

func NewHashReader(reader io.Reader, hasher hash.Hash) *HashReader {
	return &HashReader{
		inner:  reader,
		hasher: hasher,
	}
}

func (r *HashReader) Read(p []byte) (int, error) {
	n, err := r.inner.Read(p)
	if err != nil {
		return n, err
	}

	k, err := r.hasher.Write(p[:n])
	if err != nil {
		return n, err
	}

	if k != n {
		return n, errors.New("not enough byte when write to hasher")
	}

	return n, err
}

func (r *HashReader) Sum(p []byte) []byte {
	return r.hasher.Sum(p)
}
