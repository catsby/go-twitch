package twitch

import (
	"reflect"
	"testing"
)

func TestUser_Get_Id(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetUserOutput
	// c := DefaultClient()
	record(t, "users/get_by_id", func(c *Client) {
		output, err = c.GetUser(&GetUserInput{Id: 41598188})
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedUser := User{
		Id:          41598188,
		Name:        "mewnfarez",
		DisplayName: "mewnfarez",
		Type:        "user",
		Bio:         "i like to play video games",
		Logo:        "https://static-cdn.jtvnw.net/jtv_user_pictures/mewnfarez-profile_image-2af79e1168fdde0d-300x300.jpeg",
	}

	if !reflect.DeepEqual(output.User, expectedUser) {
		t.Fatalf("Error in matching users, got: \n%#v\n\nexpected:\n%#v\n\n", output.User, expectedUser)
	}
}

// Requires the user_read scope
func TestUser_Get_self(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetUserOutput
	// c := DefaultClient()
	record(t, "users/get_by_self_token", func(c *Client) {
		output, err = c.GetUser(nil)
	})
	if err != nil {
		t.Fatal(err)
	}

	expectedUser := User{
		Id:          173365798,
		Name:        "catsbygaming",
		DisplayName: "catsbygaming",
		Type:        "user",
		Bio:         "",
		Logo:        "",
	}

	if !reflect.DeepEqual(output.User, expectedUser) {
		t.Fatalf("Error in matching users, got: \n%#v\n\nexpected:\n%#v\n\n", output.User, expectedUser)
	}
}

func TestUser_Get_UserFollows(t *testing.T) {
	t.Parallel()

	var err error

	// Get
	var output *GetUserFollowsOutput
	// c := DefaultClient()
	record(t, "users/get_follows", func(c *Client) {
		output, err = c.GetUserFollows(&GetUserFollowsInput{Id: 173365798})
	})
	if err != nil {
		t.Fatal(err)
	}

	// for _, f := range output.Follows {
	// 	fmt.Printf("\tf.Channel.Game: %s\n", f.Channel.Game)
	// }
	// users := []struct {
	// 	Name string
	// 	Game string
	// }{
	// 	{
	// 		Name: "mewnfarez",
	// 		Game: "Heroes of the Storm",
	// 	},
	// 	{
	// 		Name: "kendricswissh",
	// 		Game: "Heroes of the Storm",
	// 	},
	// 	{
	// 		Name: "followgrubby",
	// 		Game: "Heroes of the Storm",
	// 	},
	// }

	// log.Printf("output: %s", spew.Sdump(output))
}
