package fileutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// Retrieve a file via a filesystem.  If the context value `filesystem` implements
// the http.FileSystem interface, it will be used to perform the retrieval in lieu
// of the local filesystem.
func RetrieveViaFilesystem(ctx context.Context, u *url.URL) (io.ReadCloser, error) {
	ctx, _ = ctxToTimeout(ctx, 0)

	var filesystem http.FileSystem = http.Dir(`/`)
	var path = strings.TrimPrefix(u.Path, `/`)

	if fs, ok := ctx.Value(`filesystem`).(http.FileSystem); ok {
		filesystem = fs
	}

	if !strings.HasPrefix(path, `/`) {
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		} else {
			return nil, fmt.Errorf("bad path: %v", err)
		}
	}

	return filesystem.Open(path)
}
