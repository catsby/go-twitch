package twitch

import (
	"fmt"
	"strconv"
	"time"
)

// VideoSort determine how video results are sorted
type VideoSortString string

const (
	VideoSortViews VideoSortString = "views"
	// Default
	VideoSortTime VideoSortString = "time"
)

// Version represents a distinct configuration version.
type Channel struct {
	Id int `mapstructure:"_id"`

	BroadcasterLanguage string `mapstructure:"broadcaster_language"`
	Description         string `mapstructure:"description"`
	DisplayName         string `mapstructure:"display_name"`
	Followers           int    `mapstructure:"followers"`
	Game                string `mapstructure:"game"`
	HTMLURL             string `mapstructure:"url"`
	Language            string `mapstructure:"language"`
	Logo                string `mapstructure:"logo"`
	Name                string `mapstructure:"name"`
	Mature              bool   `mapstructure:"mature"`
	Status              string `mapstructure:"status"`
	Views               int    `mapstructure:"views"`

	CreatedAt *time.Time `mapstructure:"created_at"`
	UpdatedAt *time.Time `mapstructure:"updated_at"`
}

// GetChannelOutput is the output of the GetChannel function.
type GetChannelOutput struct {
	Channel *Channel `mapstructure:"channel"`
}

// GetChannelInput is the input to the GetChannel function.
type GetChannelInput struct {
	Id int
}

// GetChannel returns the full list of all versions of the given service.
// GetChannel with a nil GetChannelInput, or GetChannelInput with Id 0 will hit
// the /channel/ endpoint, which retrieves the channel based on the scoped OAuth
// token. Otherwise we hit the /channels/<id> enpoint for the given ID
// See:
//  - https://dev.twitch.tv/docs/v5/reference/channels/#get-channel
//  - https://dev.twitch.tv/docs/v5/reference/channels/#get-channel-by-id
func (c *Client) GetChannel(i *GetChannelInput) (*GetChannelOutput, error) {
	path := "/channels/"
	if i == nil || i.Id == 0 {
		path = "/channel"
	} else {
		path = fmt.Sprintf("%s%d", path, i.Id)
	}

	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var o GetChannelOutput
	if err := decodeJSON(&o.Channel, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}

// GetChannelOutput is the output of the GetChannel function.
type GetChannelFollowersOutput struct {
	Total     int     `mapstructure:"_total"`
	Followers []*User `mapstructure:"follows"`
}

// GetChannelFollowersInput is the input to the GetChannelFollowers function.
type GetChannelFollowersInput struct {
	Id int
}

// GetChannelFollowers returns the full list of users following a channel
func (c *Client) GetChannelFollowers(i *GetChannelFollowersInput) (*GetChannelFollowersOutput, error) {
	path := fmt.Sprintf("/channels/%d/follows", i.Id)
	resp, err := c.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var o GetChannelFollowersOutput
	if err := decodeJSON(&o, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}

type Video struct {
	Id            string    `mapstructure:"_id"`
	BroadcastId   int       `mapstructure:"broadcast_id"`
	BroadcastType string    `mapstructure:"broadcast_type"`
	Channel       *Channel  `mapstructure:"channel"`
	CreatedAt     time.Time `mapstructure:"created_at"`
	PublishedAt   time.Time `mapstructure:"published_at"`
	Description   string    `mapstructure:"description"`
	Game          string    `mapstructure:"game"`
	Preview       *Preview  `mapstructure:"preview"`
	Title         string    `mapstructure:"title"`
	URLString     string    `mapstructure:"url"`
	Viewable      string    `mapstructure:"viewable"`
	Views         int       `mapstructure:"views"`
}

type Preview struct {
	Large    string `mapstructure:"large"`
	Medium   string `mapstructure:"medium"`
	Small    string `mapstructure:"small"`
	Template string `mapstructure:"template"`
}

// GetChannelOutput is the output of the GetChannel function.
type GetChannelVideosOutput struct {
	Total  int      `mapstructure:"_total"`
	Videos []*Video `mapstructure:"videos"`
}

// GetChannelVideosInput is the input to the GetChannelVideos function.
type GetChannelVideosInput struct {
	Id       int
	Limit    int             `mapstructure:"limit"`
	Language string          `mapstructure:"language"`
	Sort     VideoSortString `mapstructure:"sort"`
}

// GetChannelVideos returns the full list of users following a channel
func (c *Client) GetChannelVideos(i *GetChannelVideosInput) (*GetChannelVideosOutput, error) {
	path := fmt.Sprintf("/channels/%d/videos", i.Id)
	ro := new(RequestOptions)
	if i.Limit != 0 {
		ro.Params = map[string]string{
			"limit": strconv.Itoa(i.Limit),
		}
	}

	if i.Sort != "" {
		ro.Params = map[string]string{
			"sort": string(i.Sort),
		}
	}

	resp, err := c.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var o GetChannelVideosOutput
	if err := decodeJSON(&o, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}
