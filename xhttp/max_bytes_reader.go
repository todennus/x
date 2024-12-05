package xhttp

import (
	"io"
	"net/http"
)

type maxBytesReader struct {
	r   io.Reader // underlying reader
	i   int64     // max bytes initially, for MaxBytesError
	n   int64     // max bytes remaining
	err error     // sticky error
}

// MaxBytesReader is similar to http.MaxBytesReader but supporting for io.Reader
// instead of io.ReadCloser. It also doesn't support writing to any Writer when
// a MaxBytesError occurs.
func MaxBytesReader(r io.Reader, n int64) *maxBytesReader {
	return &maxBytesReader{r: r, n: n, i: 0, err: nil}
}

func (l *maxBytesReader) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}

	if len(p) == 0 {
		return 0, nil
	}

	// If they asked for a 32KB byte read but only 5 bytes are
	// remaining, no need to read 32KB. 6 bytes will answer the
	// question of the whether we hit the limit or go past it.
	// 0 < len(p) < 2^63
	if int64(len(p))-1 > l.n {
		p = p[:l.n+1]
	}
	n, err = l.r.Read(p)

	if int64(n) <= l.n {
		l.n -= int64(n)
		l.err = err
		return n, err
	}

	n = int(l.n)
	l.n = 0

	l.err = &http.MaxBytesError{Limit: l.i}
	return n, l.err
}
