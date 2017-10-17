package twitch

import (
	"encoding/json"
	"io"

	"github.com/mitchellh/mapstructure"
)

// RequestOptions is the list of options to pass to the request.
type RequestOptions struct {
	// Params is a map of key-value pairs that will be added to the Request.
	Params map[string]string

	// Headers is a map of key-value pairs that will be added to the Request.
	Headers map[string]string

	// Body is an io.Reader object that will be streamed or uploaded with the
	// Request. BodyLength is the final size of the Body.
	Body       io.Reader
	BodyLength int64
}

// decodeJSON is used to decode an HTTP response body into an interface as JSON.
func DecodeJSON(out interface{}, body io.ReadCloser) error {
	defer body.Close()

	var parsed interface{}
	dec := json.NewDecoder(body)
	if err := dec.Decode(&parsed); err != nil {
		return err
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapToHTTPHeaderHookFunc(),
			stringToTimeHookFunc(),
		),
		WeaklyTypedInput: true,
		Result:           out,
	})
	if err != nil {
		return err
	}

	// log.Printf("\nWhat is parsed: \n\n%s\n---\n", spew.Sdump(parsed))
	// log.Printf("\nWhat is out: \n\n%s\n---\n", spew.Sdump(out))
	return decoder.Decode(parsed.(map[string]interface{}))
}
