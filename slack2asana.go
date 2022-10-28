package main

import (
	"log"
	"time"
)

func main() {
	ac := NewAsanaClient()
	sc := NewSlackClient()

	err := Poll(ac, sc)
	if err != nil {
		log.Printf("%s", err)
	}

	tick := time.NewTicker(60 * time.Second)

	for {
		<-tick.C

		Poll(ac, sc)
	}
}

func Poll(ac *AsanaClient, sc *SlackClient) error {
	stars, err := sc.GetStars()
	if err != nil {
		return err
	}

	for _, item := range stars {
		if item.Type != "message" {
			continue
		}

		user, err := sc.GetUser(item.Message.User)
		if err != nil {
			return err
		}

		channel, err := sc.GetChannel(item.Channel)
		if err != nil {
			return err
		}

		title, err := sc.GetTitle(item, user, channel)
		if err != nil {
			return err
		}

		log.Printf("%s\n", title)

		err = ac.CreateTask(title)
		if err != nil {
			return err
		}

		err = sc.RemoveStar(item)
		if err != nil {
			return err
		}
	}

	return nil
}
