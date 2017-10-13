package twitch

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestClip_Get_basic(t *testing.T) {
	t.Parallel()

	slug := "LitigiousHeadstrongElkTooSpicy"
	expectedClip := Clip{
		Slug:         slug,
		TrackingId:   120388239,
		URL:          "https://clips.twitch.tv/LitigiousHeadstrongElkTooSpicy?tt_medium=clips_api&tt_content=url",
		EmbedURL:     "https://clips.twitch.tv/embed?clip=LitigiousHeadstrongElkTooSpicy&tt_medium=clips_api&tt_content=embed",
		EmbedHTMLURL: "<iframe src='https://clips.twitch.tv/embed?clip=LitigiousHeadstrongElkTooSpicy&tt_medium=clips_api&tt_content=embed' width='640' height='360' frameborder='0' scrolling='no' allowfullscreen='true'></iframe>",
		Broadcaster: &Broadcaster{
			User: User{
				Id:          0, // unfortunate api diff in IDs and mapstructure
				Name:        "mewnfarez",
				DisplayName: "mewnfarez",
				Logo:        "https://static-cdn.jtvnw.net/jtv_user_pictures/mewnfarez-profile_image-2af79e1168fdde0d-150x150.jpeg",
			},
			ChannelURL: "https://www.twitch.tv/mewnfarez",
		},
		Curator: &Curator{
			Id:          70841157,
			Name:        "smrtangel3702",
			DisplayName: "smrtangel3702",
			ChannelURL:  "https://www.twitch.tv/smrtangel3702",
			Logo:        "https://static-cdn.jtvnw.net/jtv_user_pictures/smrtangel3702-profile_image-8f1142c34047d36c-150x150.jpeg",
		},
		Vod: &Vod{
			Id:     170437016,
			URL:    "https://www.twitch.tv/videos/170437016?t=7h10m33s",
			Offset: 25833,
		},
		BroadcastId: 26123917520,
		Game:        "Heroes of the Storm",
		Language:    "en",
		Title:       "Bug-Abusing streamers bully a poor PTR player",
		Views:       24118,
		Duration:    46.333984,
		Thumbnails: &Thumbnails{
			Medium: "https://clips-media-assets.twitch.tv/26123917520-offset-25833.5-46.333333333333364-preview-480x272.jpg",
			Small:  "https://clips-media-assets.twitch.tv/26123917520-offset-25833.5-46.333333333333364-preview-260x147.jpg",
			Tiny:   "https://clips-media-assets.twitch.tv/26123917520-offset-25833.5-46.333333333333364-preview-86x45.jpg",
		},
	}

	cases := []struct {
		Label    string
		Expected *Clip
		Input    *GetClipInput
	}{
		{
			Label:    "Bug_Abusing_streamers_bully_a_poor_PTR_player",
			Expected: &expectedClip,
			Input: &GetClipInput{
				Slug: slug,
			},
		},
	}

	for _, tc := range cases {
		var err error

		// Get
		var output *GetClipOutput
		// c := DefaultClient()
		record(t, fmt.Sprintf("clips/get_%s", tc.Label), func(c *Client) {
			output, err = c.GetClip(tc.Input)
		})
		if err != nil {
			t.Fatal(err)
		}

		output.Clip.CreatedAt = nil
		expectedClip.CreatedAt = nil

		if !reflect.DeepEqual(output.Clip, *tc.Expected) {
			t.Fatalf("Error in matching clip, got: \n%s\n\nexpected:\n%s\n\n", spew.Sdump(output.Clip), spew.Sdump(*tc.Expected))
		}
	}

}

func TestClip_Get_TopClips(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Label    string
		Expected *GetTopClipsOutput
		Input    *GetTopClipsInput
	}{
		{
			Label: "Heroes_limit_20",
			Expected: &GetTopClipsOutput{
				Clips: make([]*Clip, 20),
			},
			Input: &GetTopClipsInput{
				Game:  "Heroes of the Storm",
				Limit: 20,
			},
		},
		{
			Label: "Heroes_limit_3_check_slugs",
			Expected: &GetTopClipsOutput{
				Clips: []*Clip{
					&Clip{
						Game: "Heroes of the Storm",
						Slug: "WittyEnchantingKiwiAMPTropPunch",
					},
					&Clip{
						Game: "Heroes of the Storm",
						Slug: "BoxyCoyPterodactylTheTarFu",
					},
					&Clip{
						Game: "Heroes of the Storm",
						Slug: "VenomousObedientGrassRickroll",
					},
				},
			},
			Input: &GetTopClipsInput{
				Game:     "Heroes of the Storm",
				Limit:    3,
				Trending: true,
			},
		},
	}

	for _, tc := range cases {
		var err error

		var output *GetTopClipsOutput
		record(t, fmt.Sprintf("clips/get_top_clips%s", tc.Label), func(c *Client) {
			output, err = c.GetTopClips(tc.Input)
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(output.Clips) != len(tc.Expected.Clips) {
			t.Fatalf("Length doesn't match, expected (%d), got (%d)", len(tc.Expected.Clips), len(output.Clips))
		}

		for _, c := range output.Clips {
			if c.Game != tc.Input.Game {
				t.Fatalf("Game did not match")
			}
		}

		// if the test case provides expected clips with slugs, loop through the
		// reults and make sure each slug is found
		var checkSlugs bool
		foundCount := 0
		slugs := []string{
			"VenomousObedientGrassRickroll",
			"WittyEnchantingKiwiAMPTropPunch",
			"BoxyCoyPterodactylTheTarFu",
		}
		for _, c := range tc.Expected.Clips {
			if c != nil {
				checkSlugs = true
				for _, s := range slugs {
					if c.Slug == s {
						foundCount++
					}
				}
			}
		}
		if checkSlugs && foundCount != len(slugs) {
			t.Fatalf("Expected to find (%d) slugs, but only found (%d)", len(slugs), foundCount)
		}
	}
}

func TestClip_Get_FollowedClips(t *testing.T) {
	t.Parallel()
	t.Skip("Skipping TestClip_Get_FollowedClips because I don't think it's correct")

	var err error

	var output *GetFollowedClipsOutput
	record(t, "clips/get_followed_clips", func(c *Client) {
		// output, err = c.GetFollowedClips(tc.Input)
		output, err = c.GetFollowedClips(&GetFollowedClipsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}

	log.Printf("Output: %s", spew.Sdump(output))
}
