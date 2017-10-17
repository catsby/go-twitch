package kraken

import (
	"fmt"
	"time"

	"github.com/catsby/go-twitch/twitch"
)

// User represents a user/streamer
type User struct {
	//Id          int    `mapstructure:"_id,id"` // Debug, mapstructure thing
	Id          int    `mapstructure:"_id"`
	Name        string `mapstructure:"name"`
	DisplayName string `mapstructure:"display_name"`
	Type        string `mapstructure:"type"`
	Bio         string `mapstructure:"bio"`
	Logo        string `mapstructure:"logo"`
}

// GetUserInput is the input to the GetUser function.
type GetUserInput struct {
	// Id represents a users Id, and is used for looking up users by their ID.
	// This does not require any authentication
	// Ex:
	//   GET https://api.twitch.tv/kraken/users/<user ID>
	Id int
}

// GetUsersOutput is the output of the GetUser function.
type GetUserOutput struct {
	User
}

// GetUser returns information on the user. With no GetUserInput specified, gets
// users info scoped to the access token
func (k *Kraken) GetUser(i *GetUserInput) (*GetUserOutput, error) {
	path := "/users/"
	if i == nil || i.Id == 0 {
		path = "/user"
	} else {
		path = fmt.Sprintf("%s%d", path, i.Id)
	}

	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var out GetUserOutput
	if err := twitch.DecodeJSON(&out.User, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}

// GetUserFollowsInput is the input to the GetUserFollows function.
type GetUserFollowsInput struct {
	Id int
}

type userFollows struct {
	CreatedAt     time.Time `mapstructure:"created_at"`
	Notifications bool      `mapstructure:"notifications"`
	Channel       *Channel  `mapstructure:"channel"`
}

// GetUserFollowssOutput is the output of the GetUserFollows function.
type GetUserFollowsOutput struct {
	Total   int            `mapstructure:"_total"`
	Follows []*userFollows `mapstructure:"follows"`
}

// GetUserFollows returns information on the channels a user is following.
func (k *Kraken) GetUserFollows(i *GetUserFollowsInput) (*GetUserFollowsOutput, error) {
	if i == nil || i.Id == 0 {
		return nil, fmt.Errorf("GetUserFollows requires a valid GetUserFollowsInput")
	}

	path := fmt.Sprintf("/users/%d/follows/channels", i.Id)

	resp, err := k.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var out GetUserFollowsOutput
	if err := twitch.DecodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}
