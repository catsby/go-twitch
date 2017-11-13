package helix

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/catsby/go-twitch/twitch"
	"github.com/dnaeon/go-vcr/recorder"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func TestHelixClient_DefaultClient(t *testing.T) {

	cases := []struct {
		Label       string
		Envs        []string
		Endpoint    string
		Config      *twitch.Config
		ShouldFail  bool
		ShouldError bool
	}{
		{
			Label: "Default with 1 TWITCH_ACCESS_TOKEN set",
			Envs:  []string{"TWITCH_ACCESS_TOKEN"},
		},
		{
			Label:  "Default with 1 TWITCH_ACCESS_TOKEN set, overwriting Config",
			Envs:   []string{"TWITCH_ACCESS_TOKEN"},
			Config: &twitch.Config{AccessToken: "blah"},
		},
		{
			Label:  "No env",
			Envs:   []string{},
			Config: &twitch.Config{AccessToken: "blah"},
		},
	}

	env := make(map[string]string)
	env["TWITCH_CLIENT_SECRET"] = "client_secret_123"
	env["TWITCH_ACCESS_TOKEN"] = "access_token_123"
	env["TWITCH_CLIENT_ID"] = "client_id_123"

	resetEnv := unsetEnv(t)
	defer resetEnv()
	for _, c := range cases {
		t.Run(fmt.Sprintf("%s", c.Label), func(t *testing.T) {
			unsetEnv(t)
			var accessToken string
			var clientId string
			var clientSecret string
			// var endpoint string
			if c.Config != nil {
				accessToken = c.Config.AccessToken
				clientId = c.Config.ClientId
				// endpoint = c.Config.Endpoint
			}

			for _, e := range c.Envs {
				if v, ok := env[e]; ok {
					os.Setenv(e, v)
					switch e {
					case "TWITCH_ACCESS_TOKEN":
						accessToken = env[e]
					case "TWITCH_CLIENT_SECRET":
						clientSecret = env[e]
					case "TWITCH_CLIENT_ID":
						clientId = env[e]
					}
				}
			}

			h, err := DefaultClient(c.Config)
			if err != nil {
				if !c.ShouldError {
					t.Fatalf("Error making default client when error not expected")
				}
			}

			var tErr error
			if h.accessToken != accessToken {
				tErr = fmt.Errorf("AccessToken does not match")
			}
			if h.clientId != clientId {
				tErr = fmt.Errorf("clientId does not match")
			}
			if h.clientSecret != clientSecret {
				tErr = fmt.Errorf("clientSecret does not match")
			}
			// if h.endpoint != endpoint {
			// 	t.Fatalf("Endpoint does not match")
			// }

			if c.ShouldFail && tErr == nil {
				log.Printf("should fail: %t\n", c.ShouldFail)
				log.Printf("tErr: %s\n", tErr)
				t.Fatalf("we should have failed by now")
			}
		})
	}
}

// unsetEnv unsets environment variables for testing a "clean slate" with no
// credentials in the environment
func unsetEnv(t *testing.T) func() {
	// Grab any existing Twitch keys and preserve. In some tests we'll unset these, so
	// we need to have them and restore them after
	e := getEnv()
	if err := os.Unsetenv("TWITCH_CLIENT_SECRET"); err != nil {
		t.Fatalf("Error unsetting env var TWITCH_CLIENT_SECRET: %s", err)
	}
	if err := os.Unsetenv("TWITCH_CLIENT_ID"); err != nil {
		t.Fatalf("Error unsetting env var TWITCH_CLIENT_ID: %s", err)
	}
	if err := os.Unsetenv("TWITCH_ACCESS_TOKEN"); err != nil {
		t.Fatalf("Error unsetting env var TWITCH_ACCESS_TOKEN: %s", err)
	}

	return func() {
		// re-set all the envs we unset above
		if err := os.Setenv("TWITCH_CLIENT_SECRET", e.ClientSecret); err != nil {
			t.Fatalf("Error resetting env var TWITCH_CLIENT_SECRET: %s", err)
		}
		if err := os.Setenv("TWITCH_CLIENT_ID", e.ClientId); err != nil {
			t.Fatalf("Error resetting env var TWITCH_CLIENT_ID: %s", err)
		}
		if err := os.Setenv("TWITCH_ACCESS_TOKEN", e.AccessToken); err != nil {
			t.Fatalf("Error resetting env var TWITCH_ACCESS_TOKEN: %s", err)
		}
	}
}

func setEnv(s string, t *testing.T) func() {
	e := getEnv()
	// Set all the envs to a dummy value
	if err := os.Setenv("TWITCH_CLIENT_SECRET", s); err != nil {
		t.Fatalf("Error setting env var TWITCH_CLIENT_SECRET: %s", err)
	}
	if err := os.Setenv("TWITCH_CLIENT_ID", s); err != nil {
		t.Fatalf("Error setting env var TWITCH_CLIENT_ID: %s", err)
	}
	if err := os.Setenv("TWITCH_ACCESS_TOKEN", s); err != nil {
		t.Fatalf("Error setting env var TWITCH_ACCESS_TOKEN: %s", err)
	}

	return func() {
		// re-set all the envs we unset above
		if err := os.Setenv("TWITCH_CLIENT_SECRET", e.ClientSecret); err != nil {
			t.Fatalf("Error resetting env var TWITCH_CLIENT_SECRET: %s", err)
		}
		if err := os.Setenv("TWITCH_CLIENT_ID", e.ClientId); err != nil {
			t.Fatalf("Error resetting env var TWITCH_CLIENT_ID: %s", err)
		}
		if err := os.Setenv("TWITCH_ACCESS_TOKEN", e.AccessToken); err != nil {
			t.Fatalf("Error resetting env var TWITCH_ACCESS_TOKEN: %s", err)
		}
	}
}

func getEnv() *currentEnv {
	// Grab any existing Twitch keys and preserve. In some tests we'll unset these, so
	// we need to have them and restore them after
	return &currentEnv{
		AccessToken:  os.Getenv("TWITCH_ACCESS_TOKEN"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		ClientId:     os.Getenv("TWITCH_CLIENT_ID"),
	}
}

// struct to preserve the current environment
type currentEnv struct {
	AccessToken, ClientSecret, ClientId string
}

func recordHelix(t *testing.T, fixture string, f func(*Client)) {
	modeDisabledEnv := os.Getenv("RECORD_DISABLE")
	mode := recorder.ModeReplaying
	if modeDisabledEnv == "true" {
		mode = recorder.ModeDisabled
		log.Printf("[DEBUG]: Response recording disabled")
	}

	r, err := recorder.NewAsMode("fixtures/"+fixture, mode, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := r.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	config := twitch.Config{
		HTTPClient: cleanhttp.DefaultClient(),
	}
	config.HTTPClient.Transport = r

	client, err := DefaultClient(&config)
	if err != nil {
		t.Fatal(err)
	}

	f(client)
}
