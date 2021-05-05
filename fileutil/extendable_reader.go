package fileutil

import (
	"container/list"
	"io"
	"sync"

	"github.com/ghetzel/go-stockutil/log"
)

type ExtendableReader struct {
	sources list.List
	srclock sync.Mutex
	current io.ReadCloser
}

func (self *ExtendableReader) AppendSource(rc io.ReadCloser) {
	if rc != nil {
		self.srclock.Lock()
		self.sources.PushBack(rc)
		self.srclock.Unlock()
	}
}

func (self *ExtendableReader) Read(b []byte) (int, error) {
	if err := self.closeAndAdvanceSources(); err != nil {
		return 0, err
	}

	if self.current == nil {
		return 0, io.EOF
	} else if n, err := self.current.Read(b); err == nil {
		return n, nil
	} else if err == io.EOF {
		return self.Read(b)
	} else {
		return n, err
	}
}

func (self *ExtendableReader) Close() error {
	var merr error

	self.srclock.Lock()
	defer self.srclock.Unlock()

	for e := self.sources.Front(); e != nil; e = e.Next() {
		merr = log.AppendError(merr, e.Value.(io.ReadCloser).Close())
		self.sources.Remove(e)
	}

	return merr
}

func (self *ExtendableReader) closeAndAdvanceSources() error {
	if self.current != nil {
		if err := self.current.Close(); err != nil {
			return err
		}
	}

	self.srclock.Lock()
	defer self.srclock.Unlock()

	if self.sources.Len() > 0 {
		if el := self.sources.Front(); el != nil {
			self.current = el.Value.(io.ReadCloser)
			self.sources.Remove(el)

			return nil
		}
	}

	return io.EOF
}
