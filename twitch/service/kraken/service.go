package kraken

import (
	twitch "github.com/catsby/go-twitch"
)

type Kraken struct {
	*twitch.Client
}
