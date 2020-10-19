package fileutil

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

type httpFile struct {
	*FileInfo
	buf    *bytes.Reader
	closed bool
}

func newHttpFile(data []byte) *httpFile {
	var f = &httpFile{
		FileInfo: NewFileInfo(nil),
		buf:      bytes.NewReader(data),
	}

	f.SetSize(int64(f.buf.Len()))
	f.SetIsDir(false)

	return f
}

func (self *httpFile) Close() error {
	self.closed = true
	return nil
}

func (self *httpFile) Read(b []byte) (int, error) {
	if self.closed {
		return 0, fmt.Errorf("attempted read on closed file")
	} else if self.buf == nil {
		return 0, io.EOF
	} else {
		return self.buf.Read(b)
	}
}

func (self *httpFile) Seek(offset int64, whence int) (int64, error) {
	if self.closed {
		return 0, fmt.Errorf("attempted seek on closed file")
	} else if self.buf == nil {
		return 0, io.EOF
	} else {
		return self.buf.Seek(offset, whence)
	}
}

func (self *httpFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, fmt.Errorf("not a directory")
}

func (self *httpFile) Stat() (os.FileInfo, error) {
	return self.FileInfo, nil
}
