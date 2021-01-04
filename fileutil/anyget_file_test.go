package fileutil

import (
	"io/ioutil"
	"net/url"
	"testing"

	"github.com/ghetzel/testify/assert"
)

func TestRetrieveViaFilesystem(t *testing.T) {
	var rc, rerr = RetrieveViaFilesystem(nil, &url.URL{
		Scheme: `file`,
		Path:   `testdir/a.txt`,
	})

	assert.NoError(t, rerr)

	var data, derr = ioutil.ReadAll(rc)

	assert.NoError(t, rc.Close())
	assert.NoError(t, derr)
	assert.Equal(t, "a\n", string(data))
}
