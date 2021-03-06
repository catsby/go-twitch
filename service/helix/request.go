package helix

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/catsby/go-twitch/twitch"
)

// RawRequest accepts a verb, URL, and twitch.RequestOptions struct and returns the
// constructed http.Request and any errors that occurred
func (c *Client) RawRequest(verb, p string, ro *twitch.RequestOptions) (*http.Request, error) {
	// Ensure we have request options.
	if ro == nil {
		ro = new(twitch.RequestOptions)
	}

	// Append the path to the URL.
	u := *c.url
	u.Path = strings.TrimRight(c.url.Path, "/") + "/" + strings.TrimLeft(p, "/")

	// Add the token and other params.
	var params = make(url.Values)
	for k, v := range ro.Params {
		// expand any comman seperated lists. The Helix API expects & repeated
		// parameters, not comman seperated.
		// See https://dev.twitch.tv/docs/api#requests
		s := strings.Split(v, ",")
		for _, v := range s {
			params.Add(k, v)
		}
	}

	u.RawQuery = params.Encode()

	// Create the request object.
	request, err := http.NewRequest(verb, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}

	// Set the Access Token
	if c.accessToken != "" {
		request.Header.Set(AccessTokenHeader, "Bearer "+c.accessToken)
	}

	// Set the Client Id key.
	if c.clientId != "" {
		request.Header.Set(ClientIdHeader, c.clientId)
	}

	// Set the User-Agent.
	request.Header.Set("User-Agent", twitch.UserAgent)

	// Add any custom headers.
	for k, v := range ro.Headers {
		request.Header.Add(k, v)
	}

	// Add Content-Length if we have it.
	if ro.BodyLength > 0 {
		request.ContentLength = ro.BodyLength
	}

	return request, nil
}
