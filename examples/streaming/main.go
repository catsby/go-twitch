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
}
