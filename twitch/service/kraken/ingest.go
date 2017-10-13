package kraken

// The Twitch ingesting system is the first stop for a broadcast stream. An
// ingest server receives your stream, and the ingesting system authorizes and
// registers streams, then prepares them for viewers.
type Ingest struct {
	Id           int     `mapstructure:"_id"`
	Name         string  `mapstructure:"name"`
	Default      bool    `mapstructure:"default"`
	Availability float32 `mapstructure:"availability"`
	UrlTemplate  string  `mapstructure:"url_template"`
}

// GetIngestServerListInput has no paramters
type GetIngestServerListInput struct {
}

// GetFollowedIngestsOutput is the output of the GetFollowedIngests function.
type GetIngestServerListOutput struct {
	IngestList []*Ingest `mapstructure:"ingests"`
}

// GetIngestServerList returns a list of servers for ingesting streams.
// See https://dev.twitch.tv/docs/v5/reference/ingests/#get-ingest-server-list
func (c *Client) GetIngestServerList(_ *GetIngestServerListInput) (*GetIngestServerListOutput, error) {
	resp, err := c.Get("ingests", nil)
	if err != nil {
		return nil, err
	}

	var out GetIngestServerListOutput
	if err := decodeJSON(&out, resp.Body); err != nil {
		return nil, err
	}

	return &out, nil
}
