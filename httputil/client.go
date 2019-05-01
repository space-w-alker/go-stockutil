package httputil

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ghetzel/go-stockutil/fileutil"
	"github.com/ghetzel/go-stockutil/maputil"
)

type Client struct {
	encoder         EncoderFunc
	decoder         DecoderFunc
	errorDecoder    ErrorDecoderFunc
	preRequestHook  InterceptRequestFunc
	postRequestHook InterceptResponseFunc
	uri             *url.URL
	headers         map[string]interface{}
	params          map[string]interface{}
	httpClient      *http.Client
	rootCaPool      *x509.CertPool
	insecure        bool
}

func NewClient(baseURI string) (*Client, error) {
	client := &Client{
		encoder:    JSONEncoder,
		decoder:    JSONDecoder,
		headers:    make(map[string]interface{}),
		params:     make(map[string]interface{}),
		httpClient: http.DefaultClient,
	}

	if uri, err := url.Parse(baseURI); err == nil {
		client.uri = uri
	} else {
		return nil, err
	}

	return client, nil
}

// Return a copy of the current client that uses a different encoder.
func (self *Client) WithEncoder(fn EncoderFunc) *Client {
	clientWith := new(Client)
	*clientWith = *self
	clientWith.encoder = fn
	return clientWith
}

// Return a copy of the current client that uses a different decoder.
func (self *Client) WithDecoder(fn DecoderFunc) *Client {
	clientWith := new(Client)
	*clientWith = *self
	clientWith.decoder = fn
	return clientWith
}

// Return the base URI for this client.
func (self *Client) URI() *url.URL {
	return self.uri
}

// Specify an encoder that will be used to serialize data in the request body.
func (self *Client) SetEncoder(fn EncoderFunc) {
	self.encoder = fn
}

// Specify a decoder that will be used to deserialize data in the response body.
func (self *Client) SetDecoder(fn DecoderFunc) {
	self.decoder = fn
}

// Specify a different decoder used to deserialize non 2xx/3xx HTTP responses.
func (self *Client) SetErrorDecoder(fn ErrorDecoderFunc) {
	self.errorDecoder = fn
}

// Specify a function that will be called immediately before a request is sent.
// This function has an opportunity to read and modify the outgoing request, and
// if it returns a non-nil error, the request will not be sent.
func (self *Client) SetPreRequestHook(fn InterceptRequestFunc) {
	self.preRequestHook = fn
}

// Specify a function tht will be called immediately after a response is received.
// This function is given the first opportunity to inspect the response, and if it
// returns a non-nil error, no additional processing (including the Error Decoder function)
// will be performed.
func (self *Client) SetPostRequestHook(fn InterceptResponseFunc) {
	self.postRequestHook = fn
}

// Remove all implicit HTTP request headers.
func (self *Client) ClearHeaders() {
	self.headers = nil
}

// Add an HTTP request header by name that will be included in every request. If
// value is nil, the named header will be removed instead.
func (self *Client) SetHeader(name string, value interface{}) {
	if value != nil {
		self.headers[name] = value
	} else {
		delete(self.headers, name)
	}
}

// Remove all implicit querystring parameters.
func (self *Client) ClearParams() {
	self.params = nil
}

// Add a querystring parameter by name that will be included in every request. If
// value is nil, the parameter will be removed instead.
func (self *Client) SetParam(name string, value interface{}) {
	if value != nil {
		self.params[name] = value
	} else {
		delete(self.params, name)
	}
}

// Returns the HTTP client used to perform requests
func (self *Client) Client(*http.Client) *http.Client {
	return self.httpClient
}

// Replace the default HTTP client with a user-provided one
func (self *Client) SetClient(client *http.Client) {
	self.httpClient = client
}

// Set or unset insecure TLS requests that will proceed even if the peer certificate cannot be verified.
func (self *Client) SetInsecureTLS(insecure bool) {
	defer self.syncHttpClient()
	self.insecure = insecure
}

// Set the username and password to be included in the Authorization header.
func (self *Client) SetBasicAuth(username string, password string) {
	self.SetHeader(`Authorization`, EncodeBasicAuth(username, password))
}

// Append one or more trusted certificates to the RootCA bundle that is consulted when performing HTTPS requests.
func (self *Client) AppendTrustedRootCA(pemFilenamesOrData ...interface{}) error {
	return self.updateRootCA(false, pemFilenamesOrData...)
}

// Replace the existing RootCA bundle with an explicit set of trusted certificates.
func (self *Client) SetRootCA(pemFilenamesOrData ...interface{}) error {
	return self.updateRootCA(true, pemFilenamesOrData...)
}

func (self *Client) updateRootCA(replace bool, pemFilenamesOrData ...interface{}) error {
	// when we're all done, make sure the http.Client we'll be using knows about our certs
	defer self.syncHttpClient()

	pems := make([][]byte, 0)

	for _, pem := range pemFilenamesOrData {
		if b, ok := pem.([]byte); ok {
			pems = append(pems, b)
		} else if r, ok := pem.(io.Reader); ok {
			if b, err := ioutil.ReadAll(r); err == nil {
				pems = append(pems, b)
			} else {
				return err
			}
		} else if s, ok := pem.(string); ok {
			if b, err := fileutil.ReadAll(s); err == nil {
				pems = append(pems, b)
			} else {
				return err
			}
		} else {
			return fmt.Errorf("type error: expected []byte, io.Reader, or string; got %T", pem)
		}
	}

	// clear out or setup the initial RootCA pool
	if replace {
		self.rootCaPool = x509.NewCertPool()
	} else if self.rootCaPool == nil {
		// try to inherit the RootCA from any existing client
		if transport, ok := self.httpClient.Transport.(*http.Transport); ok {
			if tcc := transport.TLSClientConfig; tcc != nil {
				self.rootCaPool = tcc.RootCAs
			}
		}

		// still nil? get the existing system CA bundle
		if self.rootCaPool == nil {
			if syspool, err := x509.SystemCertPool(); err == nil {
				self.rootCaPool = syspool
			} else {
				return fmt.Errorf("Failed to retrieve system CA pool: %v", err)
			}
		}
	}

	// append each cert to the root pool
	for _, pem := range pems {
		if !self.rootCaPool.AppendCertsFromPEM(pem) {
			return fmt.Errorf("Failed to append certificate")
		}
	}

	return nil
}

func (self *Client) syncHttpClient() {
	newTCC := &tls.Config{
		RootCAs:            self.rootCaPool,
		InsecureSkipVerify: self.insecure,
	}

	newTCC.BuildNameToCertificate()

	// handle use cases where an existing client is present
	if transport, ok := self.httpClient.Transport.(*http.Transport); ok {
		if tcc := transport.TLSClientConfig; tcc != nil {
			tcc.RootCAs = newTCC.RootCAs
			tcc.InsecureSkipVerify = self.insecure
		} else {
			transport.TLSClientConfig = newTCC
		}
	} else {
		self.httpClient.Transport = &http.Transport{
			TLSClientConfig: newTCC,
		}
	}
}

// Perform an HTTP request
func (self *Client) Request(
	method Method,
	path string,
	body interface{},
	params map[string]interface{},
	headers map[string]interface{},
) (*http.Response, error) {
	// merge given params with client-wide params
	if v, err := maputil.Merge(self.params, params); err == nil {
		params = v
	} else {
		return nil, err
	}

	// merge given headers with client-wide headers
	if v, err := maputil.Merge(self.headers, headers); err == nil {
		headers = v
	} else {
		return nil, err
	}

	reqUrl := self.uri.String()

	if path != `` {
		reqUrl = strings.Join([]string{
			strings.TrimSuffix(self.uri.String(), `/`),
			strings.TrimPrefix(path, `/`),
		}, `/`)
	}

	if url, err := url.Parse(reqUrl); err == nil {
		// set querystring values
		// ----------------------
		qs := url.Query()

		for k, v := range params {
			qs.Set(k, fmt.Sprintf("%v", v))
		}

		url.RawQuery = qs.Encode()
		// ----------------------

		if encoded, err := self.encoder(body); err == nil {
			if request, err := http.NewRequest(
				string(method),
				url.String(),
				encoded,
			); err == nil {
				// Okay, this is a little weird...
				//
				// the existing signature for EncoderFunc accepts an interface with the original
				// intention of just having a stateless transformation of one type of data into the
				// encoded version of that data.
				//
				// However, it became clear that it would be useful to ALSO allow encoders to modify
				// the request data (e.g.: add headers, set Content-Type, etc.)
				//
				// So we're going to call encoder _again_, passing it the pointer to the newly-created
				// request.  We won't check the response for errors, and just assume that if an
				// EncoderFunc implementation supports doing something with *http.Request, it will
				// Do What It Needs To Do (TM) and move on.
				//
				// This is all in service of not breaking backwards compatibility by changing this
				// function signature.  Yay.
				//
				self.encoder(request)

				// finally, because this wasn't designed right from the start and should be using
				// a context.Context, we need to special case some things we know about
				if mpfr, ok := encoded.(*multipartFormRequest); ok {
					// special case: Multipart Form requests with the boundary value in the content type
					request.Header.Set(`Content-Type`, mpfr.contentType)
				}

				for k, v := range headers {
					request.Header.Set(k, fmt.Sprintf("%v", v))
				}

				var hookObject interface{}

				if self.preRequestHook != nil {
					if v, err := self.preRequestHook(request); err == nil {
						hookObject = v
					} else {
						return nil, err
					}
				}

				// close connection after sending this request and reading its response
				request.Close = true

				// perform the request
				if response, err := self.httpClient.Do(request); err == nil {
					if self.postRequestHook != nil {
						if err := self.postRequestHook(response, hookObject); err != nil {
							return nil, err
						}
					}

					if response.StatusCode < 400 {
						return response, nil
					} else if self.errorDecoder != nil {
						return response, self.errorDecoder(response)
					} else {
						return response, fmt.Errorf("HTTP %v", response.Status)
					}
				} else {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("request init error: %v", err)
			}
		} else {
			return nil, fmt.Errorf("encoder error: %v", err)
		}
	} else {
		return nil, fmt.Errorf("url error: %v", err)
	}
}

func (self *Client) Encode(in interface{}) ([]byte, error) {
	if self.encoder != nil {
		if r, err := self.encoder(in); err == nil {
			return ioutil.ReadAll(r)
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("No encoder set")
	}
}

// Decode a response and, if applicable, automatically close the reader.
func (self *Client) Decode(r io.Reader, out interface{}) error {
	if closer, ok := r.(io.Closer); ok {
		defer closer.Close()
	}

	if self.decoder != nil {
		return self.decoder(r, out)
	} else {
		return fmt.Errorf("No decoder set")
	}
}

func (self *Client) Get(path string, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Get, path, nil, params, headers)
}

func (self *Client) GetWithBody(path string, body interface{}, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Get, path, body, params, headers)
}

func (self *Client) Post(path string, body interface{}, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Post, path, body, params, headers)
}

func (self *Client) Put(path string, body interface{}, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Put, path, body, params, headers)
}

func (self *Client) Delete(path string, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Delete, path, nil, params, headers)
}
