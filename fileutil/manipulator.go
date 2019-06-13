package fileutil

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var SkipToken = errors.New(`skip token`)

type ReadManipulatorFunc func(data []byte) ([]byte, error)

type ReadManipulator struct {
	reader    io.Reader
	fn        ReadManipulatorFunc
	splitter  bufio.SplitFunc
	scanner   *bufio.Scanner
	buffer    *bytes.Buffer
	lastToken []byte
}

func NewReadManipulator(reader io.Reader, fns ...ReadManipulatorFunc) *ReadManipulator {
	rm := &ReadManipulator{
		reader:   reader,
		splitter: bufio.ScanLines,
		buffer:   bytes.NewBuffer(nil),
	}

	if len(fns) > 0 && fns[0] != nil {
		rm.fn = fns[0]
	}

	return rm
}

func (self *ReadManipulator) Split(split bufio.SplitFunc) {
	self.splitter = split
}

func (self *ReadManipulator) Read(b []byte) (int, error) {
	if self.fn != nil {
		// initialize the scanner if we need to
		if self.scanner == nil {
			self.scanner = bufio.NewScanner(self.reader)
			self.scanner.Split(self.interceptToken)
		}

		// if there's more to scan...
		for self.scanner.Scan() {
			data := self.scanner.Bytes()

			// get the scanned bytes, run them through the manip. function
			if out, err := self.fn(data); err == nil || err == SkipToken {
				if err == nil {
					out = append(out, self.lastToken...)
				}

				self.lastToken = nil

				// write the manipulated bytes to the buffer
				if n, err := self.buffer.Write(out); err != nil {
					return n, err
				}

				// loop until we've put enough data in the buffer to satisfy the
				// requested read
				if self.buffer.Len() >= len(b) {
					break
				}
			} else {
				return 0, err
			}
		}

		// check for scan errors
		if err := self.scanner.Err(); err != nil {
			return 0, err
		}

		// return whats in the buffer, and keep doing this until its empty
		return self.buffer.Read(b)
	} else {
		return self.reader.Read(b)
	}
}

func (self *ReadManipulator) Close() error {
	if self.scanner != nil {
		self.scanner = nil
	}

	self.lastToken = nil
	self.buffer = bytes.NewBuffer(nil)

	if closer, ok := self.reader.(io.Closer); ok {
		return closer.Close()
	} else {
		return nil
	}
}

func (self *ReadManipulator) interceptToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = self.splitter(data, atEOF)

	if advance > 0 && len(data) >= advance {
		self.lastToken = append(self.lastToken, data[len(token):advance]...)
	}

	return
}
