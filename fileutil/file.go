package fileutil

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghetzel/go-stockutil/pathutil"
	isatty "github.com/mattn/go-isatty"
)

// Alias functions from pathutil as a convenience
var DirExists = pathutil.DirExists
var Exists = pathutil.Exists
var ExpandUser = pathutil.ExpandUser
var FileExists = pathutil.FileExists
var IsAppend = pathutil.IsAppend
var IsAppendable = pathutil.IsAppendable
var IsCharDevice = pathutil.IsCharDevice
var IsDevice = pathutil.IsDevice
var IsExclusive = pathutil.IsExclusive
var IsNamedPipe = pathutil.IsNamedPipe
var IsNonemptyDir = pathutil.IsNonemptyDir
var IsNonemptyExecutableFile = pathutil.IsNonemptyExecutableFile
var IsNonemptyFile = pathutil.IsNonemptyFile
var IsReadable = pathutil.IsReadable
var IsSetgid = pathutil.IsSetgid
var IsSetuid = pathutil.IsSetuid
var IsSocket = pathutil.IsSocket
var IsSticky = pathutil.IsSticky
var IsSymlink = pathutil.IsSymlink
var IsTemporary = pathutil.IsTemporary
var IsWritable = pathutil.IsWritable
var LinkExists = pathutil.LinkExists

func MustExpandUser(path string) string {
	if expanded, err := ExpandUser(path); err == nil {
		return expanded
	} else {
		panic(err.Error())
	}
}

// Return true if the given FileInfo sports a ModTime later than the current file.
func IsModifiedAfter(stat os.FileInfo, current string) bool {
	if Exists(current) {
		current = MustExpandUser(current)

		if currentStat, err := os.Stat(current); err == nil {
			// return whether the new file is modified AFTER the current one
			return stat.ModTime().After(currentStat.ModTime())
		} else {
			// fail closed; if we can't stat the existing file then don't assert that we know better
			return false
		}
	} else {
		// if the current file does not exist, then whatever file we have is newer
		return true
	}
}

func IsTerminal() bool {
	return isatty.IsTerminal(os.Stdout.Fd())
}

func ReadAll(filename string) ([]byte, error) {
	if file, err := os.Open(filename); err == nil {
		defer file.Close()
		return ioutil.ReadAll(file)
	} else {
		return nil, err
	}
}

func ReadAllString(filename string) (string, error) {
	if data, err := ReadAll(filename); err == nil {
		return string(data), nil
	} else {
		return ``, err
	}
}

func ReadAllLines(filename string) ([]string, error) {
	if data, err := ReadAllString(filename); err == nil {
		return strings.Split(data, "\n"), nil
	} else {
		return nil, err
	}
}

func ReadFirstLine(filename string) (string, error) {
	if lines, err := ReadAllLines(filename); err == nil {
		return lines[0], nil
	} else {
		return ``, err
	}
}

func MustReadAll(filename string) []byte {
	if data, err := ReadAll(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

func MustReadAllString(filename string) string {
	if data, err := ReadAllString(filename); err == nil {
		return data
	} else {
		panic(err.Error())
	}
}

// Write the contents of the given reader to the specified filename.
// Filename paths containing tilde (~) will automatically expand to the current
// user's home directory, and all intermediate parent directories will be automatically
// created.  Will return the number of bytes written, or an error.
func WriteFile(reader io.Reader, filename string) (int64, error) {
	if expanded, err := ExpandUser(filename); err == nil {
		parent := filepath.Dir(expanded)

		// create parent directory automatically
		if !DirExists(parent) {
			if err := os.MkdirAll(parent, 0700); err != nil {
				return 0, err
			}
		}

		// open the destination file for writing
		if dest, err := os.Create(expanded); err == nil {
			defer dest.Close()
			return io.Copy(dest, reader)
		} else {
			return 0, err
		}
	} else {
		return 0, err
	}
}

// Same as WriteFile, but will panic if the file cannot be written.
func MustWriteFile(reader io.Reader, filename string) int64 {
	if n, err := WriteFile(reader, filename); err == nil {
		return n
	} else {
		panic(err.Error())
	}
}
