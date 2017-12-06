package main

import (
	"fmt"
	"log"

	"github.com/catsby/go-twitch/service/kraken"
)

func main() {
	client := kraken.DefaultClient(nil)

	me, err := client.GetUser(nil)
	if err != nil {
		log.Fatalf("Error finding me: %s", err)
	}

	fmt.Println("My name is", me.Name)

	// get followed streams
	out, err := client.GetFollowedStreams(&kraken.GetFollowedStreamsInput{
		StreamType: kraken.StreamTypeLive,
	})
	if err != nil {
		log.Fatalf("Error getting followed streams: %s", err)
	}

	if len(out.Streams) == 0 {
		fmt.Println("None of your followed streams are live right now, or you have none at all")
	} else {
		fmt.Println()
		fmt.Println("Streaming now:")
	}

	for i, s := range out.Streams {
		fmt.Printf("\t%d) %s playing %s\n", i, s.Channel.DisplayName, s.Game)
	}
}
