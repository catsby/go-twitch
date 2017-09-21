package twitch

// This file was lifted and heavily modified from github.com/catsby/go-twitch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"

	"github.com/ajg/form"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
	"github.com/mitchellh/mapstructure"
)

// AccessTokenEnvVar is the name of the environment variable where the Twitch API
// key should be read from.
// ClientIdEnvVar is the name of the environment variable the client id should
// be read from.
const AccessTokenEnvVar = "TWITCH_ACCESS_TOKEN"

// Probably not needed, but offered
const ClientIdEnvVar = "TWITCH_CLIENT_ID"

// Not used yet
const ClientSecretEnvVar = "TWITCH_CLIENT_SECRET"

// AccessTokenHeader is the name of the header that contains the Twitch API key.
const AccessTokenHeader = "Authorization"
const ClientIdHeader = "Client-ID"

// DefaultEndpoint is the default endpoint for Twitch
const DefaultEndpoint = "https://api.twitch.tv/kraken/"

// ProjectURL is the url for this library.
var ProjectURL = "github.com/catsby/go-twitch"

// ProjectVersion is the version of this library.
var ProjectVersion = "0.1"

// UserAgent is the user agent for this particular client.
var UserAgent = fmt.Sprintf("catsby/go-twitch/%s (+%s; %s)",
	ProjectVersion, ProjectURL, runtime.Version())

// Client is the main entrypoint to the Twitch golang API library.
type Client struct {
	// Address is the address of Twitch's API endpoint.
	Address string

	// HTTPClient is the HTTP client to use. If one is not provided, a default
	// client will be used.
	HTTPClient *http.Client

	// accessToken is the Twitch API key to authenticate requests.
	accessToken string

	// clientId is the Twitch Application Client ID to authenticate requests.
	// Register your application here:
	//   https://dev.twitch.tv/docs/v5/guides/authentication/#registration
	// Note: probably not needed for general API consumption
	clientId string

	// url is the parsed URL from Address
	// ??
	url *url.URL
}

// DefaultClient instantiates a new Twitch API client. This function requires
// the environment variables `TWITCH_ACCESS_TOKEN` and `TWITCH_CLIENT_ID` to be
// set and contain valid access token and client id, respectively, to
// authenticate with Twitch.
func DefaultClient() *Client {
	client, err := NewClient(
		os.Getenv(ClientIdEnvVar),
		os.Getenv(AccessTokenEnvVar),
	)
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
func NewClient(clientId, accessToken string) (*Client, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("Access Token not specified")
	}

	return NewClientForEndpoint(clientId, accessToken, DefaultEndpoint)
}

// NewClientForEndpoint creates a new API client with the given key and API
// endpoint. Because Twitch allows some requests without an API key, this
// function will not error if the API token is not supplied. Attempts to make a
// request that requires an API key will return a 403 response.
func NewClientForEndpoint(clientId, accessToken string, endpoint string) (*Client, error) {
	client := &Client{accessToken: accessToken, clientId: clientId, Address: endpoint}
	return client.init()
}

func (c *Client) init() (*Client, error) {
	u, err := url.Parse(c.Address)
	if err != nil {
		return nil, err
	}
	c.url = u

	if c.HTTPClient == nil {
		c.HTTPClient = cleanhttp.DefaultClient()
	}

	return c, nil
}

func (c *Client) Get(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("GET", p, ro)
}

// Head issues an HTTP HEAD request.
func (c *Client) Head(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("HEAD", p, ro)
}

// Post issues an HTTP POST request.
func (c *Client) Post(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("POST", p, ro)
}

// PostForm issues an HTTP POST request with the given interface form-encoded.
func (c *Client) PostForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("POST", p, i, ro)
}

// Put issues an HTTP PUT request.
func (c *Client) Put(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("PUT", p, ro)
}

// PutForm issues an HTTP PUT request with the given interface form-encoded.
func (c *Client) PutForm(p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
	return c.RequestForm("PUT", p, i, ro)
}

// Delete issues an HTTP DELETE request.
func (c *Client) Delete(p string, ro *RequestOptions) (*http.Response, error) {
	return c.Request("DELETE", p, ro)
}

// Request makes an HTTP request against the HTTPClient using the given verb,
// Path, and request options.
func (c *Client) Request(verb, p string, ro *RequestOptions) (*http.Response, error) {
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
func (c *Client) RequestForm(verb, p string, i interface{}, ro *RequestOptions) (*http.Response, error) {
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

// decodeJSON is used to decode an HTTP response body into an interface as JSON.
func decodeJSON(out interface{}, body io.ReadCloser) error {
	defer body.Close()

	var parsed interface{}
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapToHTTPHeaderHookFunc(),
			stringToTimeHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           out,
	})
	if err != nil {
		return err
	}

	// log.Printf("\nWhat is parsed: \n\n%s\n---\n", spew.Sdump(parsed))
	// log.Printf("\nWhat is out: \n\n%s\n---\n", spew.Sdump(out))
	return decoder.Decode(parsed.(map[string]interface{}))
}
