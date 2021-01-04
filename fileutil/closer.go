package fileutil

import "io"

type CloserFunc = func(io.ReadCloser) error

type PostReadCloser struct {
	upstream io.ReadCloser
	closer   func(io.ReadCloser) error
}

// Implements an io.ReadCloser that can be configured to perform cleanup options whenever the
// Close() function is called.  If CloserFunc is non-nil, it will be given the upstream ReadCloser
// as an argument and will be responsible for calling Close() on it.  If nil, upstream's Close()
// function will be called directly on Close.
func NewPostReadCloser(upstream io.ReadCloser, closer CloserFunc) *PostReadCloser {
	return &PostReadCloser{
		upstream: upstream,
		closer:   closer,
	}
}

func (self *PostReadCloser) Read(b []byte) (int, error) {
	return self.upstream.Read(b)
}

func (self *PostReadCloser) Close() error {
	if fn := self.closer; fn != nil {
		return fn(self.upstream)
	} else {
		return self.upstream.Close()
	}
}
