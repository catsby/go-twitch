package kraken

import (
	"log"
	"os"
	"testing"

	"github.com/catsby/go-twitch/twitch"
	"github.com/dnaeon/go-vcr/recorder"
)

func record(t *testing.T, fixture string, f func(*twitch.Client)) {
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

	client := twitch.DefaultClient()
	client.HTTPClient.Transport = r

	f(client)
}