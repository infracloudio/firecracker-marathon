package client

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	httpCache = make(map[string]*http.Client)
	cacheLock sync.Mutex
)

// NewClient returns a new REST client.
func NewClient(host, version, userAgent string) (*Client, error) {
	baseURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	if baseURL.Path == "" {
		baseURL.Path = "localhost"
	}
	hClient := getHTTPClient(host)
	if hClient == nil {
		return nil, fmt.Errorf("Unable to parse provided url: %v", host)
	}
	c := &Client{
		base:       baseURL,
		version:    version,
		httpClient: hClient,
		userAgent:  fmt.Sprintf("%v/%v", userAgent, version),
	}
	return c, nil
}

// Client is an HTTP REST wrapper. Use one of Get/Post/Put/Delete to get a request
// object.
type Client struct {
	base       *url.URL
	version    string
	httpClient *http.Client
	userAgent  string
}

func (c *Client) SetTLS(tlsConfig *tls.Config) {
	c.httpClient = &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}
}

// Versions send a request at the /versions REST endpoint.
func (c *Client) Versions(endpoint string) ([]string, error) {
	versions := []string{}
	err := c.Get().Resource(endpoint + "/versions").Do().Unmarshal(&versions)
	return versions, err
}

// Get returns a Request object setup for GET call.
func (c *Client) Get() *Request {
	return NewRequest(c.httpClient, c.base, "GET", c.version, c.userAgent)
}

// Post returns a Request object setup for POST call.
func (c *Client) Post() *Request {
	return NewRequest(c.httpClient, c.base, "POST", c.version, c.userAgent)
}

// Put returns a Request object setup for PUT call.
func (c *Client) Put() *Request {
	return NewRequest(c.httpClient, c.base, "PUT", c.version, c.userAgent)
}

// Delete returns a Request object setup for DELETE call.
func (c *Client) Delete() *Request {
	return NewRequest(c.httpClient, c.base, "DELETE", c.version, c.userAgent)
}

func newHTTPClient(
	u *url.URL,
	tlsConfig *tls.Config,
	timeout time.Duration,
	responseTimeout time.Duration,
) *http.Client {
	httpTransport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	switch u.Scheme {
	default:
		httpTransport.Dial = func(proto, addr string) (net.Conn, error) {
			return net.DialTimeout(proto, addr, timeout)
		}
	}

	return &http.Client{Transport: httpTransport, Timeout: responseTimeout}
}

func getHTTPClient(host string) *http.Client {
	cacheLock.Lock()
	defer cacheLock.Unlock()
	c, ok := httpCache[host]
	if !ok {
		u, err := url.Parse(host)
		if err != nil {
			return nil
		}
		if u.Path == "" {
			u.Path = "/"
		}
		c = newHTTPClient(u, nil, 10*time.Second, 5*time.Minute)
		httpCache[host] = c
	}

	return c
}
