package kraken

import (
	"reflect"
	"testing"
)

func TestIngest_basic(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetIngestServerListOutput
	record(t, "ingests/get", func(c *Client) {
		output, err = c.GetIngestServerList(nil)
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedIngest := &Ingest{
		Id:           1,
		UrlTemplate:  "rtmp://live-dfw.twitch.tv/app/{stream_key}",
		Availability: 1.0,
		Name:         "US Central: Dallas, TX",
		Default:      false,
	}

	if !reflect.DeepEqual(output.IngestList[0], expectedIngest) {
		t.Fatalf("Ingest mismatch, expected (%#v), got (%#v)", expectedIngest, output.IngestList[0])
	}

	expectedCount := 44
	if len(output.IngestList) != expectedCount {
		t.Fatalf("Expected (%d) servers, got (%d)", expectedCount, len(output.IngestList))
	}
}
