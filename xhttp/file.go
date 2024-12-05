package xhttp

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/todennus/x/xbytes"
	"github.com/todennus/x/xerror"
)

const DefaultMaxFileSize = 100 * xbytes.MiB

type multipartPart struct {
	SniffReadCloser
	r *http.Request
}

func (p *multipartPart) Close() error {
	if _, err := io.Copy(io.Discard, MaxBytesReader(p.SniffReadCloser, DiscardMemory)); err != nil {
		return p.r.Body.Close()
	}

	return nil
}

type multipartFile struct {
	multipart.File
	header *multipart.FileHeader
}

type File struct {
	// Metadata
	tmpdir        string
	contentType   string
	contentLength int
	isParsed      bool
	maxFileSize   int64

	// Input
	inReader *multipartPart
	inFile   *multipartFile

	// Output
	outBytes    []byte
	outFileName string
	outFile     io.ReadSeekCloser
}

func NewFile(body SniffReadCloser, r *http.Request) *File {
	return &File{
		inReader: &multipartPart{
			SniffReadCloser: body,
			r:               r,
		},
		maxFileSize: int64(DefaultMaxFileSize),
	}
}

func NewFileFromMultipartFile(f multipart.File, header *multipart.FileHeader) *File {
	return &File{
		inFile: &multipartFile{
			File:   f,
			header: header,
		},
		maxFileSize: int64(DefaultMaxFileSize),
	}
}

func (f *File) SetMaxSize(s int64) {
	if s <= 0 {
		s = DefaultMaxFileSize
	}

	f.maxFileSize = s
}

// AsBytes converts File to bytes.Reader. It saves total file content into
// memory and suitable for small files.
func (f *File) AsBytes() (*bytes.Reader, error) {
	if f.isParsed {
		return nil, errors.New("this file was parsed before")
	}

	f.isParsed = true

	switch {
	case f.inReader != nil:
		f.outBytes = xbytes.GetMostUsedBuffer()
		buffer := bytes.NewBuffer(f.outBytes)

		if _, err := io.Copy(buffer, MaxBytesReader(f.inReader, f.maxFileSize)); err != nil {
			return nil, err
		}

		f.inReader.Close()
		f.inReader = nil
		f.outBytes = buffer.Bytes()

	case f.inFile != nil:
		if _, err := f.inFile.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		f.outBytes = xbytes.GetBuffer(f.inFile.header.Size)
		buffer := bytes.NewBuffer(f.outBytes)
		if _, err := io.Copy(buffer, MaxBytesReader(f.inFile, f.maxFileSize)); err != nil {
			return nil, err
		}

		f.inFile.Close()
		f.inFile = nil
		f.outBytes = buffer.Bytes()

	default:
		return nil, errors.New("not found any content")
	}

	return bytes.NewReader(f.outBytes), nil
}

// AsFile converts File to a ReadSeeker. The total file content is saved into a
// temporary file. It's suitable for large files.
func (f *File) AsFile() (io.ReadSeeker, error) {
	if f.isParsed {
		return nil, errors.New("this file was parsed before")
	}

	f.isParsed = true

	switch {
	case f.inFile != nil:
		if _, err := f.inFile.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}

		f.outFile = f.inFile
		f.contentLength = int(f.inFile.header.Size)
		f.inFile = nil

	case f.inReader != nil:
		copybuf := xbytes.GetBytes(32 * 1024)
		defer xbytes.PutBytes(copybuf)

		var err error
		tmpfile, err := os.CreateTemp(f.tmpdir, "xhttp-tmp-multipart-file-")
		if err != nil {
			return nil, err
		}

		f.outFileName = tmpfile.Name()

		if _, err = io.CopyBuffer(tmpfile, MaxBytesReader(f.inReader, f.maxFileSize), copybuf); err != nil {
			tmpfile.Close()
			return nil, err
		}

		tmpfile.Close() // Must close the tmpfile before opening it again for reading.
		f.inReader.Close()
		f.inReader = nil

		tmpfile, err = os.Open(f.outFileName)
		if err != nil {
			return nil, err
		}

		stat, err := tmpfile.Stat()
		if err != nil {
			return nil, err
		}

		f.outFile = tmpfile
		f.contentLength = int(stat.Size())

	default:
		return nil, errors.New("not found any content")
	}

	return f.outFile, nil
}

func (f *File) Reader() (io.Reader, error) {
	if f.isParsed {
		return nil, errors.New("this file was parsed before")
	}

	f.isParsed = true

	switch {
	case f.inReader != nil:
		return MaxBytesReader(f.inReader, f.maxFileSize), nil
	case f.inFile != nil:
		return MaxBytesReader(f.inFile, f.maxFileSize), nil
	default:
		return nil, errors.New("not found any reader")
	}
}

func (f *File) Close() error {
	var err error

	if f.inReader != nil {
		err = xerror.Join(err, f.inReader.Close())
		f.inReader = nil
	}

	if f.inFile != nil {
		err = xerror.Join(err, f.inFile.Close())
		f.inFile = nil
	}

	if f.outBytes != nil {
		xbytes.PutBuffer(f.outBytes)
		f.outBytes = nil
	}

	if f.outFile != nil {
		err = xerror.Join(err, f.outFile.Close())
		f.outFile = nil
	}

	if f.outFileName != "" {
		err = xerror.Join(err, os.Remove(f.outFileName))
		f.outFileName = ""
	}

	return err
}

func (f *File) ContentType(nsniff int64) (string, error) {
	if f.isParsed {
		return f.contentType, errors.New("this file was parsed before")
	}

	if f.contentType == "" {
		if nsniff <= 0 {
			nsniff = 512
		}

		buf := xbytes.GetBytes(nsniff)
		defer xbytes.PutBytes(buf)

		var err error

		switch {
		case f.inReader != nil:
			_, err = f.inReader.Sniff(buf)
		case f.inFile != nil:
			_, err = f.inFile.Read(buf)
			f.inFile.Seek(0, io.SeekStart)
		default:
			return "", errors.New("not found any content")
		}

		if err != nil && err != io.EOF {
			return "", err
		}

		f.contentType = http.DetectContentType(buf)
	}

	return f.contentType, nil
}

func (f *File) ContentLength() int {
	if f.contentLength <= 0 {
		switch {
		case f.inReader != nil:
			f.contentLength = 0
		case f.inFile != nil:
			f.contentLength = int(f.inFile.header.Size)
		case f.outBytes != nil:
			f.contentLength = len(f.outBytes)
		case f.outFile != nil:
			panic("content length should have been parsed in this case")
		}
	}

	return f.contentLength
}
