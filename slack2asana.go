package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type StarsResponse struct {
	Ok    bool    `json:"ok"`
	Error string  `json:"error"`
	Items []*Item `json:"items"`
}

type Item struct {
	Type    string   `json:"type"`
	Channel string   `json:"channel"`
	Message *Message `json:"message"`
}

type Message struct {
	ClientMessageId string `json:"client_msg_id"`
	Text            string `json:"text"`
	User            string `json:"user"`
	Ts              string `json:"ts"`
	Permalink       string `json:"permalink"`
}

func main() {
	c := &http.Client{}

	stars, err := getStars(c)
	if err != nil {
		panic(err)
	}

	for _, item := range stars {
		if item.Type != "message" {
			continue
		}

		fmt.Printf("%#v\n", item.Message)
	}
}

func getStars(c *http.Client) ([]*Item, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/stars.list", nil)
	if err != nil {
		return nil, err
	}

	addAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	stars := &StarsResponse{}

	err = dec.Decode(stars)
	if err != nil {
		return nil, err
	}

	if !stars.Ok {
		return nil, errors.New(stars.Error)
	}

	return stars.Items, nil
}

func addAuth(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_TOKEN")))
}
