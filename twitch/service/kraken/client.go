package kraken

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/ajg/form"
	twitch "github.com/catsby/go-twitch/twitch"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

type Client struct {
	config *twitch.Config
}

// DefaultClient instantiates a new Twitch API client. This function requires
// the environment variables `TWITCH_ACCESS_TOKEN` and `TWITCH_CLIENT_ID` to be
// set and contain valid access token and client id, respectively, to
// authenticate with Twitch.
func DefaultClient() *Client {
	// TODO not even really used right now you should fix this, it's unused
	// because I'm refactoring
	config := twitch.Config{
		ClientId:    os.Getenv(twitch.ClientIdEnvVar),
		AccessToken: os.Getenv(twitch.AccessTokenEnvVar),
		Endpoint:    twitch.DefaultEndpoint,
		HTTPClient:  nil,
	}
	client, err := NewClient(&config)
	if err != nil {
		panic(err)
	}
	return client
}

// NewClient creates a new API client with the given key and the default API
// endpoint. Twitch requires both an access token and a client id for requests,
// so we error if either is empty.
//
// Creating an access token is not yet supported by this libary
// TODO: Support creating an access token
func NewClient(config *twitch.Config) (*Client, error) {
	if config.AccessToken == "" {
		return nil, fmt.Errorf("Access Token not specified")
	}

	c := &Client{
		config: config,
	}
	return c.init()
}

func (c *Client) init() (*Client, error) {
	u, err := url.Parse(c.Endpoint)
	if err != nil {
		return nil, err
	}
	c.url = u

	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultClient()
	}

	return c, nil
}

func (c *Client) Get(p string, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.Request("GET", p, ro)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.Request("HEAD", p, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.Request("POST", p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i interface{}, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.RequestForm("POST", p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.Request("PUT", p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i interface{}, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.RequestForm("PUT", p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *twitch.RequestOptions) (*http.Response, error) {
	return c.Request("DELETE", p, ro)
}

// Request makes an HTTP request against the HTTPClient using the given verb,
// Path, and request options.
func (c *Client) Request(verb, p string, ro *twitch.RequestOptions) (*http.Response, error) {
	req, err := c.RawRequest(verb, p, ro)
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.HTTPClient.Do(req))
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// RequestForm makes an HTTP request with the given interface being encoded as
// form data.
func (c *Client) RequestForm(verb, p string, i interface{}, ro *twitch.RequestOptions) (*http.Response, error) {
	if ro == nil {
		ro = new(RequestOptions)
	}

	if ro.Headers == nil {
		ro.Headers = make(map[string]string)
	}
	ro.Headers["Content-Type"] = "application/x-www-form-urlencoded"
	// log.Printf("headers:\n%s\n", spew.Sdump(ro.Headers))

	buf := new(bytes.Buffer)
	if err := form.NewEncoder(buf).KeepZeros(true).DelimitWith('|').Encode(i); err != nil {
		return nil, err
	}
	body := buf.String()

	ro.Body = strings.NewReader(body)
	ro.BodyLength = int64(len(body))

	return c.Request(verb, p, ro)
}

// checkResp wraps an HTTP request from the default client and verifies that the
// request was successful. A non-200 request returns an error formatted to
// included any validation problems or otherwise.
func checkResp(resp *http.Response, err error) (*http.Response, error) {
	// If the err is already there, there was an error higher up the chain, so
	// just return that.
	if err != nil {
		return resp, err
	}

	switch resp.StatusCode {
	case 200, 201, 202, 204, 205, 206:
		return resp, nil
	default:
		return resp, NewHTTPError(resp)
	}
}
