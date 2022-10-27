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

type SlackClient struct {
	cli   *http.Client
	token string
}

type removeStarRequest struct {
	Channel   string `json:"channel"`
	Timestamp string `json:"timestamp"`
}

type channelResponse struct {
	Ok      bool     `json:"ok"`
	Error   string   `json:"error"`
	Channel *Channel `json:"channel"`
}

type starsResponse struct {
	Ok    bool    `json:"ok"`
	Error string  `json:"error"`
	Items []*Item `json:"items"`
}

type userResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	User  *User  `json:"user"`
}

type simpleResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
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

func NewSlackClient() *SlackClient {
	return &SlackClient{
		cli:   &http.Client{},
		token: os.Getenv("SLACK_TOKEN"),
	}
}

func (sc *SlackClient) GetStars() ([]*Item, error) {
	req, err := http.NewRequest("GET", "https://slack.com/api/stars.list", nil)
	if err != nil {
		return nil, err
	}

	sc.addAuth(req)

	resp, err := sc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	stars := &starsResponse{}

	err = dec.Decode(stars)
	if err != nil {
		return nil, err
	}

	if !stars.Ok {
		return nil, errors.New(stars.Error)
	}

	return stars.Items, nil
}

func (sc *SlackClient) GetUser(id string) (*User, error) {
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

	sc.addAuth(req)

	resp, err := sc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	user := &userResponse{}

	err = dec.Decode(user)
	if err != nil {
		return nil, err
	}

	if !user.Ok {
		return nil, errors.New(user.Error)
	}

	return user.User, nil
}

func (sc *SlackClient) GetChannel(id string) (*Channel, error) {
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

	sc.addAuth(req)

	resp, err := sc.cli.Do(req)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(resp.Body)
	channel := &channelResponse{}

	err = dec.Decode(channel)
	if err != nil {
		return nil, err
	}

	if !channel.Ok {
		return nil, errors.New(channel.Error)
	}

	return channel.Channel, nil
}

func (sc *SlackClient) RemoveStar(item *Item) error {
	body := &removeStarRequest{
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

	sc.addAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := sc.cli.Do(req)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)
	sr := &simpleResponse{}

	err = dec.Decode(sr)
	if err != nil {
		return err
	}

	if !sr.Ok {
		return errors.New(sr.Error)
	}

	return nil
}

func (sc *SlackClient) GetTitle(item *Item, user *User, channel *Channel) (string, error) {
	switch {
	case channel.IsIm:
		return fmt.Sprintf("[%s] %s", user.Name, item.Message.Text), nil
	default:
		return "", fmt.Errorf("unknown channel type: %#v", channel)
	}
}

func (sc *SlackClient) addAuth(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sc.token))
}
