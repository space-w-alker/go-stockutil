package fileutil

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type appendCoolWriter struct {
	writer io.Writer
}

func (self *appendCoolWriter) Write(b []byte) (int, error) {
	_, err := self.writer.Write([]byte(string(b) + "cool.\n"))
	return len(b), err
}

func TestDirReader(t *testing.T) {
	assert := require.New(t)

	dread := NewDirReader(`./testdir`)
	defer dread.Close()
	data, err := ioutil.ReadAll(dread)
	assert.NoError(err)
	assert.Equal("a\nb\nc\nd1\nd2\nd3\ne11\ne2\n", string(data))

	// close and see if DirReader reset properly
	assert.NoError(dread.Close())

	data, err = ioutil.ReadAll(dread)
	assert.NoError(err)
	assert.Equal("a\nb\nc\nd1\nd2\nd3\ne11\ne2\n", string(data))

	dread = NewDirReader(`./testdir/d`)
	data, err = ioutil.ReadAll(dread)
	assert.NoError(err)
	assert.Equal("d1\nd2\nd3\n", string(data))
}

func TestDirReaderSkipFunc(t *testing.T) {
	assert := require.New(t)

	dread := NewDirReader(`./testdir`)
	defer dread.Close()
	dread.SetSkipFunc(func(p string) bool {
		filename := strings.TrimSuffix(p, filepath.Ext(p))

		t.Logf("%s: %v", filename, strings.HasSuffix(filename, `1`))

		return strings.HasSuffix(filename, `1`)
	})

	data, err := ioutil.ReadAll(dread)
	assert.NoError(err)
	assert.Equal("a\nb\nc\nd2\nd3\ne2\n", string(data))
}

func TestCopyDir(t *testing.T) {
	assert := require.New(t)

	buf := bytes.NewBuffer(nil)
	cool := &appendCoolWriter{
		writer: buf,
	}

	assert.NoError(CopyDir(`./testdir`, func(path string, info os.FileInfo, err error) (io.Writer, error) {
		if info.IsDir() {
			return nil, nil
		} else {
			return cool, err
		}
	}))

	assert.EqualValues(
		"a\ncool.\nb\ncool.\nc\ncool.\nd1\ncool.\nd2\ncool.\nd3\ncool.\ne11\ncool.\ne2\ncool.\n",
		buf.String(),
	)
}
