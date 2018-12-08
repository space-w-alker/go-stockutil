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

// Loads certificates from the given file and returns a usable x509.CertPool
func LoadCertPool(filename string) (*x509.CertPool, error) {
	if data, err := fileutil.ReadAll(filename); err == nil {
		pool := x509.NewCertPool()

		if pool.AppendCertsFromPEM(data) {
			return pool, nil
		} else {
			return nil, fmt.Errorf("An error occurred adding the provided certificate(s)")
		}
	} else {
		return nil, fmt.Errorf("failed to read certificate file: %v", err)
	}
}

func updateRootCABundle(appendPem bool, client *http.Client, caBundle string) error {
	if caBundle == `` {
		return fmt.Errorf("Must specify a file to read certificates from")
	} else if client == nil {
		return fmt.Errorf("Must provide an *http.Client to modify")
	} else {
		if pool, err := LoadCertPool(caBundle); err == nil {
			if client.Transport == nil {
				client.Transport = &http.Transport{}
			}

			if htt, ok := client.Transport.(*http.Transport); ok {
				if htt.TLSClientConfig == nil {
					htt.TLSClientConfig = &tls.Config{}
				}

				htt.TLSClientConfig.RootCAs = pool
				return nil
			} else {
				return fmt.Errorf("Cannot configure TLS on HTTP transport %T", client.Transport)
			}
		} else {
			return err
		}
	}
}
