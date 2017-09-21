package twitch

import (
	"fmt"
	"reflect"
	"testing"
)

func TestChannel_Get_basic(t *testing.T) {
	t.Parallel()

	expectedSelfChannel := Channel{
		Id:                  173365798,
		BroadcasterLanguage: "",
		Description:         "",
		DisplayName:         "catsbygaming",
		Followers:           0,
		Game:                "",
		HTMLURL:             "https://www.twitch.tv/catsbygaming",
		Language:            "en",
		Logo:                "",
		Name:                "catsbygaming",
		Mature:              false,
		Status:              "",
		Views:               0,
	}

	expectedChannel := Channel{
		Id:                  8822,
		BroadcasterLanguage: "en",
		Description:         "Dedicated to creating the most epic entertainment experiences... ever.",
		DisplayName:         "Blizzard",
		Followers:           238514,
		Game:                "Gamescom 2017",
		HTMLURL:             "https://www.twitch.tv/blizzard",
		Language:            "en",
		Logo:                "https://static-cdn.jtvnw.net/jtv_user_pictures/blizzard-profile_image-ede74b5c45a43edc-300x300.jpeg",
		Name:                "blizzard",
		Mature:              false,
		Status:              "Blizzard at gamescom",
		Views:               8797750,
	}

	cases := []struct {
		Label    string
		Expected *Channel
		Input    *GetChannelInput
	}{
		{
			Label:    "nil_input",
			Expected: &expectedSelfChannel,
			Input:    nil,
		},
		{
			Label:    "zero_id",
			Expected: &expectedSelfChannel,
			Input:    &GetChannelInput{},
		},
		{
			Label:    "other",
			Expected: &expectedChannel,
			Input:    &GetChannelInput{Id: 8822},
		},
	}

	for _, tc := range cases {
		var err error

		// Get
		var output *GetChannelOutput
		//c := DefaultClient()
		record(t, fmt.Sprintf("channels/get_%s", tc.Label), func(c *Client) {
			output, err = c.GetChannel(tc.Input)
		})
		if err != nil {
			t.Fatal(err)
		}

		output.Channel.CreatedAt = nil
		output.Channel.UpdatedAt = nil

		if !reflect.DeepEqual(output.Channel, tc.Expected) {
			t.Fatalf("Error in matching channel, got: \n%#v\n\nexpected:\n%#v\n\n", output.Channel, tc.Expected)
		}
	}

}

func TestChannel_GetFollows_basic(t *testing.T) {
	t.Parallel()

	// Get
	var err error
	var output *GetChannelFollowersOutput
	// c := DefaultClient()
	record(t, "channels/follows", func(c *Client) {
		output, err = c.GetChannelFollowers(&GetChannelFollowersInput{Id: 43664778})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := &GetChannelFollowersOutput{
		Total: 27821,
	}

	if output.Total != expectedOutput.Total {
		t.Fatalf("Total count mismatch, expected (%d), got (%d)", expectedOutput.Total, output.Total)
	}

	// Didn't want to make a []*User slice with actual users, just check expected
	// count
	expectedFollowerCount := 25
	if len(output.Followers) != expectedFollowerCount {
		t.Fatalf("Total followers count mismatch, expected (%d), got (%d)", expectedFollowerCount, len(output.Followers))
	}
}

func TestChannel_Get_ChannelVideos_basic(t *testing.T) {
	t.Parallel()

	// Get
	var err error
	var output *GetChannelVideosOutput
	record(t, "channels/videos_basic", func(c *Client) {
		output, err = c.GetChannelVideos(&GetChannelVideosInput{Id: 43664778})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := &GetChannelVideosOutput{
		Total:  110,
		Videos: make([]*Video, 10, 10),
	}

	if output.Total != expectedOutput.Total {
		t.Fatalf("Total count mismatch, expected (%d), got (%d)", expectedOutput.Total, output.Total)
	}

	if len(output.Videos) != len(expectedOutput.Videos) {
		t.Fatalf("Total count mismatch, expected (%d), got (%d)", len(expectedOutput.Videos), len(output.Videos))
	}

}
