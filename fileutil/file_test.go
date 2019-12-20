package fileutil

import (
	"testing"

	"github.com/ghetzel/testify/require"
)

func TestSetExt(t *testing.T) {
	assert := require.New(t)

	assert.Equal(``, SetExt(``, ``))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, ``))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, `.jpg`, `.bmp`))
	assert.Equal(`/nothingburger.txt`, SetExt(`/nothingburger.txt`, `jpg`, `.bmp`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`, `.txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`, `.txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `.jpg`, `txt`))
	assert.Equal(`/nothingburger.jpg`, SetExt(`/nothingburger.txt`, `jpg`, `txt`))

	assert.Equal(`/nothingburger.info.xml`, SetExt(`/nothingburger.info.json`, `xml`))
	assert.Equal(`/nothingburger.info.xml`, SetExt(`/nothingburger.info.json`, `.xml`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `.xml`, `info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `.xml`, `.info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `xml`, `info.json`))
	assert.Equal(`/nothingburger.xml`, SetExt(`/nothingburger.info.json`, `xml`, `.info.json`))
}
