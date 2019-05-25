package fileutil

import (
	"os"
	"time"
)

// An os.FileInfo-compatible wrapper that allows for individual values to be overridden.
type FileInfo struct {
	os.FileInfo
	name  string
	size  int64
	mode  *os.FileMode
	mtime time.Time
	dir   *bool
	sys   interface{}
}

func NewFileInfo(wrap ...os.FileInfo) *FileInfo {
	if len(wrap) > 0 && wrap[0] != nil {
		return &FileInfo{
			FileInfo: wrap[0],
		}
	} else {
		return new(FileInfo)
	}
}

func (self *FileInfo) Name() string {
	if self.name != `` {
		return self.name
	} else if self.FileInfo != nil {
		return self.FileInfo.Name()
	} else {
		return ``
	}
}

func (self *FileInfo) Size() int64 {
	if self.size != 0 {
		return self.size
	} else if self.FileInfo != nil {
		return self.FileInfo.Size()
	} else {
		return 0
	}
}

func (self *FileInfo) Mode() os.FileMode {
	if self.mode != nil {
		return *self.mode
	} else if self.FileInfo != nil {
		return self.FileInfo.Mode()
	} else {
		return 0
	}
}

func (self *FileInfo) ModTime() time.Time {
	if !self.mtime.IsZero() {
		return self.mtime
	} else if self.FileInfo != nil {
		return self.FileInfo.ModTime()
	} else {
		return time.Time{}
	}
}

func (self *FileInfo) IsDir() bool {
	if self.dir != nil {
		return *self.dir
	} else if self.FileInfo != nil {
		return self.FileInfo.IsDir()
	} else {
		return false
	}
}

func (self *FileInfo) Sys() interface{} {
	if self.sys != nil {
		return self.sys
	} else if self.FileInfo != nil {
		return self.FileInfo.Sys()
	} else {
		return nil
	}
}

func (self *FileInfo) SetName(name string) {
	self.name = name
}

func (self *FileInfo) SetSize(sz int64) {
	self.size = sz
}

func (self *FileInfo) SetMode(mode os.FileMode) {
	self.mode = &mode
}

func (self *FileInfo) SetModTime(mtime time.Time) {
	self.mtime = mtime
}

func (self *FileInfo) SetIsDir(isDir bool) {
	self.dir = &isDir
}

func (self *FileInfo) SetSys(iface interface{}) {
	self.sys = iface
}
