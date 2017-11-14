package twitch

// This file was lifted and heavily modified from github.com/catsby/go-twitch

import (
	"fmt"
	"net/http"
	"runtime"
)

// Default (v5) vs. "New" API: Twitch is rolling out a new API and deprecating
// v5. We include a client for both v5 (kraken) jand the new API (helix), with
// Helix endpoints and features being added as they come. The Fields below apply
// to both clients unless specified.
//
// See:
// - https://dev.twitch.tv/docs/api

// DefaultEndpoint is the default endpoint for Twitch
const KrakenEndpoint = "https://api.twitch.tv/kraken/"
const HelixEndpoint = "https://api.twitch.tv/helix/"
const DefaultEndpoint = KrakenEndpoint

// ProjectURL is the url for this library.
var ProjectURL = "github.com/catsby/go-twitch"

// ProjectVersion is the version of this library.
var ProjectVersion = "0.1"

// UserAgent is the user agent for this particular client.
var UserAgent = fmt.Sprintf("catsby/go-twitch/%s (+%s; %s)",
	ProjectVersion, ProjectURL, runtime.Version())

// Holds the configuraiton options for all clients.
type Config struct {
	// Endpoint is the address of Twitch's API endpoint.
	Endpoint string

	// HTTPClient is the HTTP client to use. If one is not provided, a default
	// client will be used.
	HTTPClient *http.Client

	// accessToken is the Twitch API key to authenticate requests.
	AccessToken string

	// clientId is the Twitch Application Client ID to authenticate requests.
	// Register your application here:
	//   https://dev.twitch.tv/docs/v5/guides/authentication/#registration
	// Note: probably not needed for general API consumption
	ClientId string

	// this is probably used?
	ClientSecret string
}

// DefaultClient instantiates a new Twitch API client for talking to the new
// Helix API. This function requires the environment variables
// `TWITCH_ACCESS_TOKEN` and `TWITCH_CLIENT_ID` to be set and contain valid
// access token and client id, respectively, to authenticate with Twitch.
// func HelixClient() *Client {
// 	client, err := NewClient(
// 		os.Getenv(ClientIdEnvVar),
// 		os.Getenv(AccessTokenEnvVar),
// 		HelixEndpoint,
// 		nil,
// 	)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return client
// }
