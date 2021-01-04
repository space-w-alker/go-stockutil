package fileutil

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ghetzel/go-stockutil/stringutil"
	"github.com/ghetzel/go-stockutil/typeutil"
	"github.com/ghetzel/testify/assert"
)

func testHttpServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var _, act, arg = stringutil.SplitTriple(req.URL.Path, `/`)

		switch act {
		case `sleep`:
			time.Sleep(typeutil.Duration(arg))
		}

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("OK"))
	}))
}

func TestRetrieve(t *testing.T) {
	var server = testHttpServer(t)
	defer server.Close()

	for uri, expected := range map[string]string{
		"testdir/a.txt":         "a\n",
		"./testdir/a.txt":       "a\n",
		"file:///testdir/a.txt": "a\n",
		server.URL + "/hello":   "OK",
	} {
		var rc, err = Retrieve(nil, uri)
		assert.NoError(t, err, uri)

		var actual, rerr = ioutil.ReadAll(rc)

		assert.NoError(t, rc.Close(), uri)
		assert.NoError(t, rerr, uri)
		assert.Equal(t, expected, string(actual), uri)
	}
}
