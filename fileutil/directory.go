package fileutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type DirReaderOption int

const (
	NoRecursive DirReaderOption = iota
	FailOnError
)

type DirReaderOptions []DirReaderOption

func (self DirReaderOptions) Has(option DirReaderOption) bool {
	for _, opt := range self {
		if opt == option {
			return true
		}
	}

	return false
}

type SkipFunc func(string) bool

// A DirReader provides a streaming io.Reader interface to all files in a given
// directory, with options for handling unreadable entries and recursion.
type DirReader struct {
	root         string
	options      DirReaderOptions
	loaded       bool
	entries      []os.FileInfo
	size         int64
	currentEntry int
	current      io.ReadCloser
	skipFn       SkipFunc
}

func NewDirReader(path string, options ...DirReaderOption) *DirReader {
	return &DirReader{
		root:    path,
		options: DirReaderOptions(options),
	}
}

// Set a function that will be called for each path encountered while reading.
// If this function returns true, that path (and its descedants) will not be read.
func (self *DirReader) SetSkipFunc(fn SkipFunc) {
	self.skipFn = fn
}

func (self *DirReader) setup() error {
	if rt, entries, err := self.readDir(); err == nil {
		self.root = rt
		self.current = nil
		self.currentEntry = 0
		self.entries = entries
		self.size = 0
		self.loaded = true
		return nil
	} else {
		return err
	}
}

func (self *DirReader) advanceAndRead(b []byte) (int, error) {
	if self.current != nil {
		if err := self.current.Close(); err != nil {
			return 0, err
		}
	}

	// if the current entry is in-bounds
	if self.currentEntry < len(self.entries) {
		entry := self.entries[self.currentEntry]
		self.currentEntry += 1
		path := filepath.Join(self.root, entry.Name())

		// if a skipFunc we provided, and it returned false, return a skiperr from here
		if skipFn := self.skipFn; skipFn != nil && skipFn(path) {
			return self.advanceAndRead(b)
		}

		if entry.IsDir() {
			if !self.options.Has(NoRecursive) {
				// directories that we recurse into are just instances of DirReaders whose root is the directory
				subreader := NewDirReader(path, self.options...)
				subreader.SetSkipFunc(self.skipFn)

				self.current = subreader
				return self.current.Read(b)
			} else {
				return self.advanceAndRead(b)
			}
		} else if file, err := os.Open(path); err == nil {
			self.current = file
			return self.current.Read(b)
		} else {
			return 0, err
		}
	} else {
		return 0, io.EOF
	}
}

func (self *DirReader) Read(b []byte) (int, error) {
	// do initial setup if we're starting out
	if !self.loaded {
		if err := self.setup(); err != nil {
			return 0, err
		}
	}

	// check if we have a current file
	if self.current != nil {
		if n, err := self.current.Read(b); err == nil {
			// if so, read from that file
			return n, nil
		} else if err == io.EOF || !self.options.Has(FailOnError) {
			// if the current file is EOF (or we're skipping errors), advance to the next one and start reading it
			// keep advancing until the error is not skipEntryErr
			return self.advanceAndRead(b)
		} else {
			return n, err
		}
	} else {
		return self.advanceAndRead(b)
	}
}

// close open files and reset the internal reader
func (self *DirReader) Close() error {
	if self.current != nil {
		self.current.Close()
	}

	self.loaded = false

	return nil
}

// read the immediate entries from  the current root directory, or if the current root
// is a file, treat it like a directory of one entry
func (self *DirReader) readDir() (string, []os.FileInfo, error) {
	if root, err := ExpandUser(self.root); err == nil {
		if DirExists(root) {
			entries, err := ioutil.ReadDir(self.root)

			sort.Slice(entries, func(i int, j int) bool {
				return entries[i].Name() < entries[j].Name()
			})

			return root, entries, err
		} else if stat, err := os.Stat(root); err == nil {
			return root, []os.FileInfo{stat}, nil
		} else {
			return ``, nil, err
		}
	} else {
		return ``, nil, err
	}
}

type CopyEntryFunc func(path string, info os.FileInfo, err error) (io.Writer, error)

// Recursively walk the entries of a given directory, calling CopyEntryFunc for each entry.
// The io.Writer returned from the function will have that file's contents written to it.  If
// the io.Writer is nil, the file will not be written anywhere but no error will be returned.
// If CopyEntryFunc returns an error, the behavior will be consistent with filepath.WalkFunc
func CopyDir(root string, fn CopyEntryFunc) error {
	if p, err := ExpandUser(root); err == nil {
		root = p
	} else {
		return err
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if w, err := fn(path, info, err); err == nil && w != nil {
			if file, err := os.Open(path); err == nil {
				defer file.Close()
				_, err := io.Copy(w, file)
				return err
			} else {
				return err
			}
		} else if err != nil {
			return err
		} else {
			return nil
		}
	})
}
