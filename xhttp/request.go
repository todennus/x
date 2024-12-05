package xhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/todennus/x/mime"
	"github.com/todennus/x/xbytes"
	"github.com/todennus/x/xreflect"
)

type FormDataRequest interface {
	NumFiles() int
}

const maxBufferBytes = 512

const DefaultMaxInMemoryMultipartSize = 10 * xbytes.MiB // 10MiB
var MaxInMemoryMultipartSize = DefaultMaxInMemoryMultipartSize
var DiscardMemory = 512 * xbytes.KiB

func ParseHTTPRequest[T any](req *http.Request) (*T, error) {
	return ParseHTTPRequestWithLimitedSize[T](req, -1)
}

func ParseHTTPRequestWithLimitedSize[T any](req *http.Request, limitedSize int64) (*T, error) {
	if limitedSize > 0 {
		req.Body = http.MaxBytesReader(nil, req.Body, limitedSize)
	}

	var t T

	if err := parseURLParameter(&t, req); err != nil {
		return nil, translateError(err)
	}

	if err := parseURLQuery(&t, req); err != nil {
		return nil, translateError(err)
	}

	switch req.Method {
	case http.MethodGet:
		return &t, nil

	case http.MethodPost, http.MethodPut, http.MethodDelete:
		contentType := req.Header.Get("Content-Type")

		switch {
		case contentType == mime.ApplicationJSON:
			if err := parseJSONBody(&t, req); err != nil {
				return nil, translateError(err)
			}
			return &t, nil

		case contentType == mime.ApplicationFormURLEncoded:
			if err := parseURLEncodedFormData(&t, req); err != nil {
				return nil, translateError(err)
			}

			return &t, nil

		case strings.HasPrefix(contentType, mime.MultipartFormData):
			usingReader := false
			if freq, ok := (any)(&t).(FormDataRequest); ok {
				if freq.NumFiles() == 1 {
					usingReader = true
				}
			}

			var err error
			if usingReader {
				err = parseMultipartFormdataReader(&t, req)
			} else {
				err = parseMultipartFormData(&t, req)
			}

			if err != nil {
				return nil, translateError(err)
			}

			return &t, nil

		default:
			if contentType == "" {
				contentType = "<empty>"
			}

			return nil, fmt.Errorf("%wnot support content type %s", ErrHTTPBadRequest, contentType)
		}

	default:
		return nil, fmt.Errorf("%wnot support method %s", ErrHTTPBadRequest, req.Method)
	}
}

func parseURLQuery(obj any, req *http.Request) error {
	query := req.URL.Query()
	return parse(obj, req, false, "query", func(r *http.Request, s string) any {
		if len(query[s]) == 0 {
			return ""
		}

		return strings.Join(query[s], " ")
	})
}

func parseURLParameter(obj any, req *http.Request) error {
	return parse(obj, req, false, "param", func(r *http.Request, s string) any {
		return chi.URLParam(r, s)
	})
}

func parseJSONBody(obj any, req *http.Request) error {
	m := map[string]any{}
	if err := json.NewDecoder(req.Body).Decode(&m); err != nil {
		return fmt.Errorf("%winvalid json", ErrHTTPBadRequest)
	}

	return parse(obj, req, true, "json", func(r *http.Request, s string) any {
		s, _, _ = strings.Cut(s, ",")
		return m[s]
	})
}

func parseURLEncodedFormData(obj any, req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}

	return parse(obj, req, false, "form", func(req *http.Request, fieldName string) any {
		if len(req.Form[fieldName]) == 0 {
			return ""
		}

		return strings.Join(req.Form[fieldName], " ")
	})
}

func parseMultipartFormData(obj any, req *http.Request) error {
	if err := req.ParseMultipartForm(MaxInMemoryMultipartSize); err != nil {
		return err
	}

	files := map[string]multipart.File{}
	headers := map[string]*multipart.FileHeader{}

	return parse(obj, req, false, "multipart", func(r *http.Request, s string) any {
		key, tag, _ := strings.Cut(s, ",")
		if tag == "file" {
			if _, ok := files[key]; !ok {
				file, header, err := r.FormFile(key)
				if err != nil {
					if !errors.Is(err, http.ErrMissingFile) {
						slog.Warn("failed-to-part-form-file", "source", "xhttp/request.go", "err", err)
					}

					return nil
				}

				files[key] = file
				headers[key] = header
			}

			return NewFileFromMultipartFile(files[key], headers[key])
		}

		return strings.Join(r.PostForm[key], " ")
	})
}

func parseMultipartFormdataReader(obj any, req *http.Request) error {
	mreader, err := req.MultipartReader()
	if err != nil {
		return err
	}

	var file *multipart.Part
	values := map[string]string{}
	for {
		part, err := mreader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		key := part.FormName()
		if part.FileName() != "" {
			file = part
			break
		} else {
			err := func() error {
				buf := xbytes.GetBytes(maxBufferBytes)
				defer xbytes.PutBytes(buf)

				for {
					n, err := part.Read(buf)
					if n > 0 {
						values[key] += string(buf[:n])
					}

					if err == io.EOF {
						break
					}

					if err != nil {
						return err
					}
				}

				return nil
			}()

			if err != nil {
				return err
			}
		}
	}

	return parse(obj, req, false, "multipart", func(r *http.Request, s string) any {
		key, tag, _ := strings.Cut(s, ",")
		if tag == "file" {
			if file == nil {
				return nil
			}

			return NewFile(NewSniffReadCloser(file), r)
		}

		return values[key]
	})
}

func parse(obj any, req *http.Request, strict bool, tagName string, fieldVal func(*http.Request, string) any) error {
	return xreflect.Parse(obj, strict, tagName, func(s string) any {
		return fieldVal(req, s)
	})
}

func translateError(err error) error {
	if mberr := AsMaxBytesError(err); mberr != nil {
		return fmt.Errorf("%wtoo large request (limit %d bytes)", ErrHTTPTooLarge, mberr.Limit)
	}

	if errors.Is(err, xreflect.ErrBadFormat) {
		return fmt.Errorf("%w%s", ErrHTTPBadRequest, err.Error())
	}

	return err
}
