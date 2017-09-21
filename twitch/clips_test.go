package twitch

import (
	"fmt"
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
		Views:       24109,
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
