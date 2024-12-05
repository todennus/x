package xhttp

import (
	"errors"
	"net/http"
)

var (
	ErrHTTPBadRequest = errors.New("")
	ErrHTTPTooLarge   = errors.New("")
)

func AsMaxBytesError(err error) *http.MaxBytesError {
	var mberr *http.MaxBytesError
	if errors.As(err, &mberr) {
		return mberr
	}

	return nil
}
