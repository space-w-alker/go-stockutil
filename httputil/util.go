package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/ghetzel/go-stockutil/fileutil"
)

// Configures the given http.Client to accept TLS certificates validated by the given PEM-encoded CA bundle file
func SetRootCABundle(client *http.Client, caBundle string) error {
	return updateRootCABundle(false, client, caBundle)
}

// Configures the given http.Client to accept TLS certificates validated by the given PEM-encoded CA bundle file,
// appending the certificates to any existing CA bundle.
func AppendRootCABundle(client *http.Client, caBundle string) error {
	return updateRootCABundle(true, client, caBundle)
}

func updateRootCABundle(appendPem bool, client *http.Client, caBundle string) error {
	if caBundle != `` && client != nil {
		if data, err := fileutil.ReadAll(caBundle); err == nil {
			if client.Transport == nil {
				client.Transport = &http.Transport{}
			}

			if htt, ok := client.Transport.(*http.Transport); ok {
				if htt.TLSClientConfig == nil {
					htt.TLSClientConfig = &tls.Config{}
				}

				if htt.TLSClientConfig.RootCAs == nil || !appendPem {
					htt.TLSClientConfig.RootCAs = x509.NewCertPool()
				}

				if htt.TLSClientConfig.RootCAs.AppendCertsFromPEM(data) {
					return nil
				} else {
					return fmt.Errorf("An error occurred configuring TLS root CAs")
				}
			} else {
				return fmt.Errorf("Cannot configure TLS on HTTP transport %T", client.Transport)
			}
		} else {
			return fmt.Errorf("failed to read CA bundle: %v", err)
		}
	} else {
		return nil
	}
}
