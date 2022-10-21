package main

import (
	"net/http"
)

type StarsResponse struct {
	Ok bool `json:"ok"`
	Items []Item `json:"items"`
}

type Item struct {
	Type string `json:"type"`
	Channel string `json:"channel",omitempty`
	Message *Message `json:"message"`
}

type Message Struct {
	ClientMessageId string `json:"client_message_id"`
	Text string `json:"text"`
	User string `json:"user"`
	Ts string `json:"ts"`
	Permalink string `json:"permalink"`
}


