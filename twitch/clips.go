package twitch

import (
	"fmt"
	"time"
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
func (c *Client) GetClip(i *GetClipInput) (*GetClipOutput, error) {
	if i == nil || i.Slug == "" {
		return nil, fmt.Errorf("[ERR] No Slug for GetClip")
	}
	path := fmt.Sprintf("/%s/%s", "clips", i.Slug)

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var o GetClipOutput
	if err := decodeJSON(&o.Clip, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}
