package httputil

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/ghetzel/go-stockutil/maputil"
)

type Client struct {
	uri        *url.URL
	headers    map[string]interface{}
	params     map[string]interface{}
	httpClient *http.Client
}

func NewClient(baseURI string) (*Client, error) {
	client := &Client{
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

func (self *Client) AddHeader(name string, value interface{}) {
	self.headers[name] = value
}

func (self *Client) AddParam(name string, value interface{}) {
	self.params[name] = value
}

func (self *Client) SetClient(client *http.Client) {
	self.httpClient = client
}

func (self *Client) Request(
	method Method,
	path string,
	body io.Reader,
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

	if url, err := url.Parse(
		fmt.Sprintf(
			"%v/%v",
			strings.TrimSuffix(self.uri.String(), `/`),
			strings.TrimPrefix(path, `/`),
		),
	); err == nil {
		// set querystring values
		// ----------------------
		qs := url.Query()

		for k, v := range params {
			qs.Set(k, fmt.Sprintf("%v", v))
		}

		url.RawQuery = qs.Encode()
		// ----------------------

		if request, err := http.NewRequest(
			string(method),
			url.String(),
			body,
		); err == nil {
			for k, v := range headers {
				request.Header.Set(k, fmt.Sprintf("%v", v))
			}

			// perform the request
			return self.httpClient.Do(request)
		} else {
			return nil, fmt.Errorf("request init error: %v", err)
		}
	} else {
		return nil, fmt.Errorf("url error: %v", err)
	}
}

func (self *Client) Get(path string, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Get, path, nil, params, headers)
}

func (self *Client) GetWithBody(path string, body io.Reader, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Get, path, body, params, headers)
}

func (self *Client) Post(path string, body io.Reader, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Post, path, body, params, headers)
}

func (self *Client) Put(path string, body io.Reader, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Put, path, body, params, headers)
}

func (self *Client) Delete(path string, params map[string]interface{}, headers map[string]interface{}) (*http.Response, error) {
	return self.Request(Delete, path, nil, params, headers)
}
