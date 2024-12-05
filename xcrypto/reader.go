package xcrypto

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"hash"
	"io"

	"github.com/todennus/x/xbytes"
)

func Sha256(reader io.Reader) (string, error) {
	return Hash(sha256.New(), reader)
}

func Hash(hasher hash.Hash, reader io.Reader) (string, error) {
	return HashByChunk(hasher, reader, 32*xbytes.KiB)
}

func HashByChunk(hasher hash.Hash, content io.Reader, chunk int64) (string, error) {
	if chunk <= 0 {
		chunk = 1024
	}

	buffer := xbytes.GetBytes(chunk)
	defer xbytes.PutBytes(buffer)

	for {
		n, err := content.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}

		hasher.Write(buffer[:n])

		if err == io.EOF {
			break
		}
	}

	return base64.RawURLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

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
