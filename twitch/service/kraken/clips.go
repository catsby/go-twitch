package kraken

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/catsby/go-twitch/twitch"
)

// Version represents a distinct configuration version.
type Clip struct {
	Slug string

	// Embeded types
	Broadcaster *Broadcaster `mapstructure:"broadcaster"`
	Curator     *Curator     `mapstructure:"curator"`
	Vod         *Vod         `mapstructure:"vod"`
	Thumbnails  *Thumbnails  `mapstructure:"thumbnails"`

	TrackingId   int    `mapstructure:"tracking_id"`
	URL          string `mapstructure:"url"`
	EmbedURL     string `mapstructure:"embed_url"`
	EmbedHTMLURL string `mapstructure:"embed_html"`

	BroadcastId int     `mapstructure:"broadcast_id"`
	Game        string  `mapstructure:"game"`
	Language    string  `mapstructure:"language"`
	Title       string  `mapstructure:"title"`
	Views       int     `mapstructure:"views"`
	Duration    float64 `mapstructure:"duration"`

	CreatedAt *time.Time `mapstructure:"created_at"`
}

type Thumbnails struct {
	Medium string `mapstructure:"medium"`
	Small  string `mapstructure:"small"`
	Tiny   string `mapstructure:"tiny"`
}

type Curator struct {
	Id          int    `mapstructure:"id"`
	Name        string `mapstructure:"name"`
	DisplayName string `mapstructure:"display_name"`
	ChannelURL  string `mapstructure:"channel_url"`
	Logo        string `mapstructure:"logo"`
}

type Vod struct {
	Id     int    `mapstructure:"id"`
	URL    string `mapstructure:"url"`
	Offset int    `mapstructure:"offset"`
}

type Broadcaster struct {
	// The Broadcaster is a user/streamer, unfortunately, the API returns the
	// broadcater ID without the normal _id prefix that User has (uses just id
	// here) , so mapstructure doesn't parse it out. The result is we get User
	// fields as we expect, except ID is 0
	//
	// This could be addressed with a method on Broadcaster that looks up the id
	// by the name.
	User       `mapstructure:",squash"`
	ChannelURL string `mapstructure:"channel_url"`
}

// GetClipOutput is the output of the GetClip function.
type GetClipOutput struct {
	Clip
}

// GetClipInput is the input to the GetClip function.
type GetClipInput struct {
	// Clips are referenced by a globally unique string called a slug
	Slug string
}

// Gets details about a specified clip
// See:
//  - https://dev.twitch.tv/docs/v5/reference/clips#get-clip
func (k *Client) GetClip(i *GetClipInput) (*GetClipOutput, error) {
	if i == nil || i.Slug == "" {
		return nil, fmt.Errorf("[ERR] No Slug for GetClip")
	}
	path := fmt.Sprintf("/%s/%s", "clips", i.Slug)

	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var o GetClipOutput
	if err := twitch.DecodeJSON(&o.Clip, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}

// GetTopClipsOutput is the output of the GetTopClips function.
type GetTopClipsOutput struct {
	Clips  []*Clip `mapstructure:"clips"`
	Cursor string  `mapstructure:"_cursor"`
}

// GetTopClipsInput is the input to the GetTopClips function.
type GetTopClipsInput struct {
	// Channel name. If this is specified, top clips for only this channel are returned; otherwise, top clips for all channels are returned. If both channel and game are specified, game is ignored.
	Channel string
	// Tells the server where to start fetching the next set of results, in a multi-page response.
	Cursor string
	// Game name. (Game names can be retrieved with the Search Games endpoint.) If this is specified, top clips for only this game are returned; otherwise, top clips for all games are returned. If both channel and game are specified, game is ignored.
	Game string
	//   Comma-separated list of languages, which constrains the languages of videos returned. Examples: es, en,es,th. If no language is specified, all languages are returned. Default: "". Maximum: 28 languages.
	Language string
	// Maximum number of most-recent objects to return. Default: 10. Maximum: 100.
	Limit int
	// The window of time to search for clips. Valid values: day, week, month, all. Default: week.
	Period string
	//   If true, the clips returned are ordered by popularity; otherwise, by viewcount. Default: false.
	Trending bool
}

// Gets top clips
// See:
//  - https://dev.twitch.tv/docs/v5/reference/clips#get-top-clips
func (k *Client) GetTopClips(i *GetTopClipsInput) (*GetTopClipsOutput, error) {
	path := fmt.Sprintf("/clips/top")
	params := map[string]string{
		"channel":  i.Channel,
		"cursor":   i.Cursor,
		"game":     i.Game,
		"language": i.Language,
		"limit":    strconv.Itoa(i.Limit),
		"period":   i.Period,
		"trending": strconv.FormatBool(i.Trending),
	}

	// TODO should refactor into a method that returns a twitch.RequestOptions struct for
	// a given GetTopClipsInput
	// strip out any empty parameters
	for k, v := range params {
		if v == "" {
			delete(params, k)
		}
		// need limit, default is 10, if it's omitted, just remove it
		if k == "limit" && v == "0" {
			delete(params, k)
		}
	}

	ro := twitch.RequestOptions{
		Params: params,
	}

	resp, err := k.Get(path, &ro)

	if err != nil {
		return nil, err
	}

	var o GetTopClipsOutput
	if err := twitch.DecodeJSON(&o, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}

////////
////////
////////
////////

type pageable struct {
	// Tells the server where to start fetching the next set of results, in a
	// multi-page response.
	Cursor string `mapstructure:"cursor"`

	// Maximum number of most-recent objects to return. Default: 10. Maximum: 100.
	Limit int `mapstructure:"limit"`

	// The window of time to search for clips. Valid values: day, week, month,
	// all. Default: week.
	Period string `mapstructure:"period"`
}

// GetFollowedClipsOutput is the output of the GetFollowedClips function.
type GetFollowedClipsOutput struct {
	*pageable

	Clips []*Clip `mapstructure:"clips"`
}

// GetFollowedClipsInput is the input to the GetFollowedClips function.
type GetFollowedClipsInput struct {
	*pageable

	// Game name. (Game names can be retrieved with the Search Games endpoint.) If this is specified, top clips for only this game are returned; otherwise, top clips for all games are returned. If both channel and game are specified, game is ignored.
	Game string
	//   Comma-separated list of languages, which constrains the languages of videos returned. Examples: es, en,es,th. If no language is specified, all languages are returned. Default: "". Maximum: 28 languages.
	Language string
	//   If true, the clips returned are ordered by popularity; otherwise, by viewcount. Default: false.
	Trending bool
}

// Gets Followed clips
// TODO/Note: I'm not sure how this api is supposed to return things, I'm
// getting [] for a reply and I can't seem to "follow" clips right now. I don't
// consider this endpoint of the SDK to "work"
// Scope: user_read
// See:
//  - https://dev.twitch.tv/docs/v5/reference/clips#get-top-clips
func (k *Client) GetFollowedClips(i *GetFollowedClipsInput) (*GetFollowedClipsOutput, error) {
	log.Printf("[WARN] GetFollowedClips probably doesn't acually work")
	resp, err := k.Get("/clips/followed", nil)

	if err != nil {
		return nil, err
	}

	var o GetFollowedClipsOutput
	if err := twitch.DecodeJSON(&o, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}
