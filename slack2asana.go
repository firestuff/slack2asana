package main

import (
	"log"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	ac := NewAsanaClient()
	sc := NewSlackClient()

	err := Poll(ac, sc)
	if err != nil {
		log.Printf("%s", err)
	}

	for {
		time.Sleep(time.Duration(rand.Intn(60)) * time.Second)

		err = Poll(ac, sc)
		if err != nil {
			log.Printf("%s", err)
		}
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

		title, err := sc.GetTrimmedTitle(item, user, channel)
		if err != nil {
			return err
		}

		notes, err := sc.GetNotes(item, user, channel)
		if err != nil {
			return err
		}

		log.Printf("%s\n", title)

		err = ac.CreateTask(title, notes)
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
