package helix

import (
	"fmt"
	"testing"
)

func TestGame_Get_basic(t *testing.T) {
	t.Parallel()

	expectedGameHeroes := Game{
		Id:   "32959",
		Name: "Heroes of the Storm",
	}
	expectedGameFortnite := Game{
		Id:   "33214",
		Name: "Fortnite",
	}

	cases := []struct {
		Label    string
		Expected []*Game
		Input    *GetGamesInput
	}{
		{
			Label:    "SingleGame",
			Expected: []*Game{&expectedGameHeroes},
			Input: &GetGamesInput{
				Names: []string{"Heroes of the Storm"},
			},
		},
		{
			Label:    "DoubleGame",
			Expected: []*Game{&expectedGameHeroes, &expectedGameFortnite},
			Input: &GetGamesInput{
				Names: []string{"Heroes of the Storm", "Fortnite"},
			},
		},
		{
			Label:    "DoubleGameReversed",
			Expected: []*Game{&expectedGameHeroes, &expectedGameFortnite},
			Input: &GetGamesInput{
				Names: []string{"Fortnite", "Heroes of the Storm"},
			},
		},
		{
			Label:    "ByIds",
			Expected: []*Game{&expectedGameHeroes, &expectedGameFortnite},
			Input: &GetGamesInput{
				Ids: []int{32959, 33214},
			},
		},
		{
			Label:    "NameAndId",
			Expected: []*Game{&expectedGameHeroes, &expectedGameFortnite},
			Input: &GetGamesInput{
				Names: []string{"Heroes of the Storm"},
				Ids:   []int{33214},
			},
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Helix Get Games - %s", tc.Label), func(t *testing.T) {
			var err error

			// Get
			var output *GetGamesOutput

			recordHelix(t, fmt.Sprintf("games/get_%s", tc.Label), func(c *Client) {
				output, err = c.GetGames(tc.Input)
			})
			if err != nil {
				t.Fatal(err)
			}

			expectedCount := len(tc.Expected)
			if len(output.Games) != expectedCount {
				t.Fatalf("Expected at least (%d) game(s) here, but got (%d)", expectedCount, len(output.Games))
			}

			games := make(map[string]string)
			for _, g := range tc.Expected {
				games[g.Id] = g.Name
			}

			for _, g := range output.Games {
				if f, ok := games[g.Id]; ok {
					if f != g.Name {
						t.Fatalf("Name mismatch: %s / %s", f, g.Name)
					}
					delete(games, g.Id)
				}
			}

			if len(games) > 0 {
				t.Fatalf("Failed to match all games, not found: %#v", games)
			}
		})
	}
}
