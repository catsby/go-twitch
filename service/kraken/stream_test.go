package kraken

import (
	"log"
	"net/http"
	"testing"

	"github.com/catsby/go-twitch/twitch"
)

func TestStream_GetFollowed(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetFollowedStreamsOutput
	recordKraken(t, "streams/followed", func(c *kraken.Client) {
		output, err = c.GetFollowedStreams(&GetFollowedStreamsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := 1
	if output.Total != expectedOutput {
		t.Fatalf("Expectec output.Total to be (%d), got (%d)", expectedOutput, output.Total)
	}

	for _, s := range output.Streams {
		if s.Channel == nil {
			t.Fatalf("A livestream channel should not be nil")
		}

		if s.Game != "Heroes of the Storm" {
			t.Fatalf("Expected heroes")
		}
	}
}

func TestStream_GetStream_basic(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetStreamOutput
	var resp *http.Response
	recordKraken(t, "streams/41598188", func(c *kraken.Client) {
		output, resp, err = c.GetStream(&GetStreamInput{
			ChannelId: 41598188,
		})
	})
	if err != nil {
		log.Printf("HTTP Response: %#v", resp)
		t.Fatal(err)
	}

	if output.Stream == nil {
		t.Fatalf("Found nil Stream")
	}

	if output.Stream.Game != "Heroes of the Storm" {
		t.Fatalf("Expected 'Heroes of the Storm' game, got (%s)", output.Stream.Game)
	}
}

func TestStream_GetLiveStreams(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Path        string
		Req         *GetLiveStreamsInput
		StreamCount int
		TotalCount  int
	}{
		{
			Path:        "streams/live",
			Req:         &GetLiveStreamsInput{},
			StreamCount: 25,
			TotalCount:  22940,
		},
	}

	for i, c := range cases {
		var err error
		// Get
		var output *GetLiveStreamsOutput
		recordKraken(t, c.Path, func(client *Client) {
			output, err = client.GetLiveStreams(c.Req)
		})
		if err != nil {
			t.Fatal(err)
		}

		// 25 in fixture at time
		if len(output.Streams) != c.StreamCount {
			t.Fatalf("Case (%d) expected (%d) live streams, got (%d)", i, c.StreamCount, len(output.Streams))
		}

		// Total was 19944
		if output.Total != c.TotalCount {
			t.Fatalf("Case (%d) error in total count, expected (%d), got (%d)", i, c.TotalCount, output.Total)
		}
	}
}

func TestStream_GetStream_Summary(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetStreamSummaryOutput
	recordKraken(t, "stream/summary/overwatch_summary", func(c *kraken.Client) {
		output, err = c.GetStreamSummary(&GetStreamSummaryInput{
			Game: "Overwatch",
		})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedChannels := 864
	if output.Channels != expectedChannels {
		t.Fatalf("Expectec output.Channels to be (%d), got (%d)", expectedChannels, output.Channels)
	}

	expectedViewers := 30316
	if output.Viewers != expectedViewers {
		t.Fatalf("Expectec output.Viewers to be (%d), got (%d)", expectedViewers, output.Viewers)
	}

}

func TestStream_GetFeatured(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetFeaturedStreamsOutput
	recordKraken(t, "streams/featured", func(c *kraken.Client) {
		output, err = c.GetFeaturedStreams(&GetFeaturedStreamsInput{})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedTotal := 25
	if output.Total != expectedTotal {
		t.Fatalf("Expectec output.Total to be (%d), got: (%d)", expectedTotal, output.Total)
	}

	featuredStream := output.FeaturedStreams[0]

	expectedGame := "Dead by Daylight"
	if featuredStream.Stream.Game != expectedGame {
		t.Fatalf("Expected game to be (%s), got: (%s)", expectedGame, featuredStream.Stream.Game)
	}

	if featuredStream.Priority != 5 {
		t.Fatalf("Expected stream priority to be 5, got (%d)", featuredStream.Priority)
	}
}
