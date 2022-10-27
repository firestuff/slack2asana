package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type RemoveStarRequest struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
}

type ChannelResponse struct {
	Ok      bool     `json:"ok"`
	Error   string   `json:"error"`
	Channel *Channel `json:"channel"`
}

type StarsResponse struct {
	Ok    bool    `json:"ok"`
	Error string  `json:"error"`
	Items []*Item `json:"items"`
}

type UserResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	User  *User  `json:"user"`
}

type Channel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	IsChannel bool   `json:"is_channel"`
	IsGroup   bool   `json:"is_group"`
	IsIm      bool   `json:"is_im"`
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

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SimpleResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
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

		user, err := getUser(c, item.Message.User)
		if err != nil {
			panic(err)
		}

		channel, err := getChannel(c, item.Channel)
		if err != nil {
			panic(err)
		}

		title, err := getTitle(item, user, channel)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", title)

		err = removeStar(c, item)
		if err != nil {
			panic(err)
		}
	}
}

func getTitle(item *Item, user *User, channel *Channel) (string, error) {
	switch {
	case channel.IsIm:
		return fmt.Sprintf("[%s] %s", user.Name, item.Message.Text), nil
	default:
		return "", fmt.Errorf("unknown channel type: %#v", channel)
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

func getUser(c *http.Client, id string) (*User, error) {
	u, err := url.Parse("https://slack.com/api/users.info")
	if err != nil {
		return nil, err
	}

	v := &url.Values{}
	v.Add("user", id)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	addAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	user := &UserResponse{}

	err = dec.Decode(user)
	if err != nil {
		return nil, err
	}

	if !user.Ok {
		return nil, errors.New(user.Error)
	}

	return user.User, nil
}

func getChannel(c *http.Client, id string) (*Channel, error) {
	u, err := url.Parse("https://slack.com/api/conversations.info")
	if err != nil {
		return nil, err
	}

	v := &url.Values{}
	v.Add("channel", id)
	u.RawQuery = v.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	addAuth(req)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	channel := &ChannelResponse{}

	err = dec.Decode(channel)
	if err != nil {
		return nil, err
	}

	if !channel.Ok {
		return nil, errors.New(channel.Error)
	}

	return channel.Channel, nil
}

func removeStar(c *http.Client, item *Item) error {
	body := &RemoveStarRequest{
		Channel:   item.Channel,
		Timestamp: item.Message.Ts,
	}

	js, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/stars.remove", bytes.NewReader(js))
	if err != nil {
		return err
	}

	addAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	sr := &SimpleResponse{}

	err = dec.Decode(sr)
	if err != nil {
		return err
	}

	if !sr.Ok {
		return errors.New(sr.Error)
	}

	return nil
}

func addAuth(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("SLACK_TOKEN")))
}
