package xhttp

import (
	"io"
	"net/http"

	"github.com/todennus/x/mime"
)

type Sniffer interface {
	// Sniff is similar to Read but does not change the cursor.
	Sniff([]byte) (int, error)
}

type SniffReadCloser interface {
	io.ReadCloser
	Sniffer
}

// DetectContentType uses http.DetectContentType for a Sniffer instead.
func DetectContentType(r Sniffer, n int) string {
	p := make([]byte, n)
	_, err := r.Sniff(p)
	if err != nil && err != io.EOF {
		return mime.ApplicationOctetStream
	}

	return http.DetectContentType(p)
}

type sniffReadSeekCloser struct {
	io.ReadSeekCloser
}

func (r *sniffReadSeekCloser) Sniff(p []byte) (int, error) {
	defer r.ReadSeekCloser.Seek(0, io.SeekStart)

	n, err := r.ReadSeekCloser.Read(p)
	if err != nil {
		return n, err
	}

	return n, nil
}

func NewSniffReadSeekCloser(r io.ReadSeekCloser) SniffReadCloser {
	return &sniffReadSeekCloser{ReadSeekCloser: r}
}

type sniffReadCloser struct {
	sniff []byte
	inner io.ReadCloser
}

func NewSniffReadCloser(r io.ReadCloser) SniffReadCloser {
	if r == nil {
		return nil
	}

	return &sniffReadCloser{sniff: nil, inner: r}
}

func (r *sniffReadCloser) Sniff(p []byte) (int, error) {
	n, err := r.inner.Read(p)
	r.sniff = make([]byte, n)
	copy(r.sniff, p[:n])

	return n, err
}

func (r *sniffReadCloser) Read(p []byte) (int, error) {
	nSniff := len(r.sniff)
	if nSniff > len(p) {
		nSniff = len(p)
	}

	if nSniff > 0 {
		copy(p, r.sniff[:nSniff])
		r.sniff = r.sniff[nSniff:]
	}

	if nSniff == len(p) {
		return nSniff, nil
	}

	n, err := r.inner.Read(p[nSniff:])
	return nSniff + n, err
}

func (r *sniffReadCloser) Close() error {
	r.sniff = nil
	return r.inner.Close()
}
