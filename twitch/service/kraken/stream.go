package kraken

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/catsby/go-twitch/twitch"
)

// StreamType is the type of cache action.
type StreamType string

const (
	// StreamTypeCache sets the cache to cache.
	StreamTypeLive StreamType = "live"

	// StreamTypePass sets the cache to pass through.
	StreamTypePlayList StreamType = "playlist"

	// StreamTypeRestart sets the cache to restart the request.
	StreamTypeAll StreamType = "all"
)

// Stream represents a distinct stream
type Stream struct {
	Id         int      `mapstructure:"_id"` // maybe _id ?
	AverageFps int      `mapstructure:"average_fps"`
	Delay      int      `mapstructure:"delay"`
	Game       string   `mapstructure:"game"`
	Channel    *Channel `mapstructure:"channel"`
}

// GetFollowedStreamsInput is the input to the GetFollowedStreams function.
type GetFollowedStreamsInput struct {
	// Maximum number of objects to return. Default: 25. Maximum: 100.O
	Limit int `mapstructure:"limit"`

	// Constrains the type of streams returned. Valid values: live, playlist, all.
	// Playlists are offline streams of VODs (Video on Demand) that appear live.
	// Default: live.
	StreamType StreamType `mapstructure:"stream_type"`

	//Object offset for pagination of results. Default: 0.
	Offset int `mapstructure:"offset"`
}

// GetFollowedStreamsOutput is the output of the GetFollowedStreams function.
type GetFollowedStreamsOutput struct {
	// Total results returned
	Total int `mapstructure:"_total"`

	// List of matching Streams
	Streams []*Stream
}

// GetFollowedStreams returns a list of online streams a user is following,
// based on a specified OAuth token.
func (k *Kraken) GetFollowedStreams(i *GetFollowedStreamsInput) (*GetFollowedStreamsOutput, error) {
	path := fmt.Sprintf("/streams/followed")
	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var out GetFollowedStreamsOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}

// GetStreamInput is the input to the GetStream function.
type GetStreamInput struct {
	ChannelId int
	// Constrains the type of streams returned. Valid values: live, playlist, all. Playlists are offline streams of VODs (Video on Demand) that appear live. Default: live.
	StreamType StreamType
}

// GetStreamOutput is the output of the GetStream function.
type GetStreamOutput struct {
	// List of matching Streams
	Stream *Stream
}

// GetStream returns the full list of all versions of the given service.
func (k *Kraken) GetStream(i *GetStreamInput) (*GetStreamOutput, *http.Response, error) {
	if i == nil || i.ChannelId == 0 {
		return nil, nil, errors.New("Invalid GetStreamInput: ChannelId is required and cannot be zero")
	}
	path := fmt.Sprintf("/streams/%d", i.ChannelId)
	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, resp, err
	}

	var out GetStreamOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, resp, err
	}

	return &out, resp, nil
}

// GetLiveStreamsInput is the input to the GetStream function.
type GetLiveStreamsInput struct {
	// Slice of channel ids to filter
	ChannelIds []string `mapstructure:"channel"`

	// Game name to filter on
	Game string `mapstructure:"game"`

	// Constrains streams to a language.
	// Ex: en, fi, es-mx. Only one language can be specified.
	// Omit to search default "all languages"
	Language string `mapstructure:"language"`

	// Maximum number of objects to return. Default: 25. Maximum: 100.O
	Limit int `mapstructure:"limit"`

	// Constrains the type of streams returned. Valid values: live, playlist, all. Playlists are offline streams of VODs (Video on Demand) that appear live. Default: live.
	StreamType StreamType `mapstructure:"stream_type"`

	//Object offset for pagination of results. Default: 0.
	Offset int `mapstructure:"offset"`
}

// GetLiveStreamsOutput is the output of the GetStream function.
type GetLiveStreamsOutput struct {
	// Total
	Total int `mapstructure:"_total"`

	// List of matching Streams
	Streams []*Stream
}

// GetStream returns the full list of all versions of the given service.
func (k *Kraken) GetLiveStreams(i *GetLiveStreamsInput) (*GetLiveStreamsOutput, error) {
	path := "/streams"
	ro := new(twitch.RequestOptions)
	// for _,s:=range i.
	if i.Game != "" {
		ro.Params = map[string]string{
			"game": i.Game,
		}
	}
	resp, err := k.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var out GetLiveStreamsOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}

// GetStreamSummaryInput is the input to the GetStreamSummary function.
type GetStreamSummaryInput struct {
	// Game name to filter on
	Game string `mapstructure:"game"`
}

// GetLiveStreamsOutput is the output of the GetStream function.
type GetStreamSummaryOutput struct {
	Channels int `mapstructure:"channels"`
	Viewers  int `mapstructure:"viewers"`
}

// GetStream returns the full list of all versions of the given service.
func (k *Kraken) GetStreamSummary(i *GetStreamSummaryInput) (*GetStreamSummaryOutput, error) {
	path := "/streams/summary"
	ro := new(twitch.RequestOptions)
	// for _,s:=range i.
	if i.Game != "" {
		ro.Params = map[string]string{
			"game": i.Game,
		}
	}
	resp, err := k.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var out GetStreamSummaryOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}

type FeaturedStream struct {
	Image     string
	Priority  int
	Scheduled bool
	Sponsored bool
	Stream    *Stream
}

// GetFeaturedStreamsInput is the input to the GetFeaturedStreams function.
type GetFeaturedStreamsInput struct {
	// Maximum number of objects to return. Default: 25. Maximum: 100.O
	Limit int `mapstructure:"limit"`

	// Constrains the type of streams returned. Valid values: live, playlist, all. Playlists are offline streams of VODs (Video on Demand) that appear live. Default: live.
	StreamType StreamType `mapstructure:"stream_type"`

	//Object offset for pagination of results. Default: 0.
	Offset int `mapstructure:"offset"`
}

// GetFeaturedStreamsOutput is the output of the GetFeaturedStreams function.
type GetFeaturedStreamsOutput struct {
	// Total results returned
	Total int

	// List of matching Streams
	FeaturedStreams []*FeaturedStream `mapstructure:"featured"`
}

// GetFeaturedStreams returns the full list of all versions of the given service.
func (k *Kraken) GetFeaturedStreams(i *GetFeaturedStreamsInput) (*GetFeaturedStreamsOutput, error) {
	path := fmt.Sprintf("/streams/featured")
	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var out GetFeaturedStreamsOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	out.Total = len(out.FeaturedStreams)

	return &out, nil
}
