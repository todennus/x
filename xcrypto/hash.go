package xcrypto

import (
	"crypto/sha256"
	"errors"
	"hash"
	"io"

	"github.com/todennus/x/xbytes"
)

func Sha256(reader io.Reader) ([]byte, error) {
	return Hash(sha256.New(), reader)
}

func Hash(hasher hash.Hash, reader io.Reader) ([]byte, error) {
	return HashByChunk(hasher, reader, 32*xbytes.KiB)
}

func HashByChunk(hasher hash.Hash, content io.Reader, chunk int64) ([]byte, error) {
	if chunk <= 0 {
		chunk = 1024
	}

	buffer := xbytes.GetBytes(chunk)
	defer xbytes.PutBytes(buffer)

	for {
		n, err := content.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		if n > 0 {
			hasher.Write(buffer[:n])
		}

		if err == io.EOF {
			break
		}
	}

	return hasher.Sum(nil), nil
}
