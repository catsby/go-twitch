package kraken

import (
	twitch "github.com/catsby/go-twitch/twitch"
)

type Client struct {
	*twitch.Client
}
