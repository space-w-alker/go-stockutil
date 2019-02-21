package fileutil

import "io"

type PostReadCloser struct {
	upstream io.ReadCloser
	closer   func(io.ReadCloser) error
}

// Implements an io.ReadCloser that can be configured to perform cleanup options whenever the
// Close() function is called.
func NewPostReadCloser(upstream io.ReadCloser, closer func(io.ReadCloser) error) *PostReadCloser {
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
		return nil
	}
}
