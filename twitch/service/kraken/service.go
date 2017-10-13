package kraken

import (
	twitch "github.com/catsby/go-twitch/twitch"
)

type Kraken struct {
	*twitch.Client
}
