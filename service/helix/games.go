package helix

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/catsby/go-twitch/twitch"
)

// Version represents a distinct configuration version.
type Game struct {
	Id   string `mapstructure:"id"`
	Name string `mapstructure:"name"`
}

// GetGamesOutput is the output of the GetGames function.
type GetGamesOutput struct {
	Games []*Game `mapstructure:"data"`
}

// GetGamesInput is the input to the GetGames function.
type GetGamesInput struct {
	// Games are referenced by a globally unique string called a slug
	Names []string `mapstructure:"name"`
	Ids   []int    `mapstructure:"id"`
}

// Gets details specific games. Can be list of games, ids
// See:
//  - https://dev.twitch.tv/docs/api/reference#get-games
func (k *Client) GetGames(i *GetGamesInput) (*GetGamesOutput, error) {
	if i == nil || (len(i.Names) == 0 && len(i.Ids) == 0) {
		return nil, fmt.Errorf("[ERR] No Name or Id for GetGamess")
	}
	path := "/games"

	ro := new(twitch.RequestOptions)
	if len(i.Names) > 0 {
		ro.Params = map[string]string{
			"name": strings.Join(i.Names, ","),
		}
	}

	if len(i.Ids) > 0 {
		var ids []string
		if ro.Params == nil {
			ro.Params = map[string]string{}
		}
		for _, i := range i.Ids {
			ids = append(ids, strconv.Itoa(i))
		}
		ro.Params["id"] = strings.Join(ids, ",")
	}

	resp, err := k.Get(path, ro)
	if err != nil {
		return nil, err
	}

	var o GetGamesOutput
	if err := twitch.DecodeJSON(&o, resp.Body); err != nil {
		return nil, err
	}

	return &o, nil
}
