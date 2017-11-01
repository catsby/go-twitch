package kraken

import (
	"log"
	"os"
	"testing"

	"github.com/catsby/go-twitch/twitch"
	"github.com/dnaeon/go-vcr/recorder"
	cleanhttp "github.com/hashicorp/go-cleanhttp"
)

func recordKraken(t *testing.T, fixture string, f func(*Client)) {
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

	client := DefaultClient(&config)

	f(client)
}
