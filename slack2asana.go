package main

import (
	"fmt"
)

func main() {
	sc := NewSlackClient()

	stars, err := sc.GetStars()
	if err != nil {
		panic(err)
	}

	for _, item := range stars {
		if item.Type != "message" {
			continue
		}

		user, err := sc.GetUser(item.Message.User)
		if err != nil {
			panic(err)
		}

		channel, err := sc.GetChannel(item.Channel)
		if err != nil {
			panic(err)
		}

		title, err := sc.GetTitle(item, user, channel)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", title)

		err = sc.RemoveStar(item)
		if err != nil {
			panic(err)
		}
	}
}
