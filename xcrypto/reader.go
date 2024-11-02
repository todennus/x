package xcrypto

import (
	"errors"
	"hash"
	"io"
)

var _ io.Reader = (*HashReader)(nil)

type HashReader struct {
	reader io.Reader
	hasher hash.Hash
}

// NewHashReader returns a special Reader in which an inner hasher consumes the
// data when the Read() method is called. Use the Sum() method to retrieve the
// hash result.
func NewHashReader(reader io.Reader, hasher hash.Hash) *HashReader {
	return &HashReader{
		reader: reader,
		hasher: hasher,
	}
}

func (r *HashReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
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

func (r *HashReader) Sum() []byte {
	return r.hasher.Sum(nil)
}
