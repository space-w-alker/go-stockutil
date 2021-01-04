package fileutil

import (
	"context"
	"io/ioutil"
	"net/url"
	"testing"
	"time"

	"github.com/ghetzel/testify/assert"
)

func TestRetrieveViaHTTP(t *testing.T) {
	var server = testHttpServer(t)
	defer server.Close()

	var ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	var u, err = url.Parse(server.URL)

	assert.NoError(t, err)

	var rc, rerr = RetrieveViaHTTP(ctx, u)
	assert.NoError(t, rerr)

	var data, derr = ioutil.ReadAll(rc)

	assert.NoError(t, rc.Close())
	assert.NoError(t, derr)
	assert.Equal(t, "OK", string(data))

	u.Path = `/sleep/1100ms`

	rc, rerr = RetrieveViaHTTP(ctx, u)
	assert.Contains(t, rerr.Error(), `context deadline exceeded`)
}
