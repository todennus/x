package xhttp

import (
	"errors"
	"net/url"
)

func ParseURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("invalid scheme")
	}

	if u.User != nil {
		return nil, errors.New("not support userinfo in url")
	}

	return u, nil
}
