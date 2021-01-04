package fileutil

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/ghetzel/go-stockutil/typeutil"
)

var HttpDefaultTimeout = 10 * time.Second

// Retrieve a file via HTTP or HTTPS.
func RetrieveViaHTTP(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	ctx, timeout := ctxToTimeout(ctx, HttpDefaultTimeout)

	var client = &http.Client{
		Timeout: timeout,
	}

	if typeutil.Bool(ctx.Value(`insecure`)) {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if res, err := client.Get(u.String()); err == nil {
		if res.StatusCode < 300 {
			if res.Body != nil {
				return res.Body, nil
			} else {
				return NullReadCloser, nil
			}
		} else {
			return nil, fmt.Errorf("responded HTTP: %v", res.Status)
		}
	} else {
		return nil, err
	}
}
