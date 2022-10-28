package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type AsanaClient struct {
	cli       *http.Client
	token     string
	workspace string
	assignee  string
	project   string
}

type addTaskRequest struct {
	Data *addTaskRequestInt `json:"data"`
}

type addTaskRequestInt struct {
	Name      string   `json:"name"`
	Workspace string   `json:"workspace"`
	Assignee  string   `json:"assignee"`
	Projects  []string `json:"projects"`
}

func NewAsanaClient() *AsanaClient {
	return &AsanaClient{
		cli:       &http.Client{},
		token:     os.Getenv("ASANA_TOKEN"),
		workspace: os.Getenv("ASANA_WORKSPACE"),
		assignee:  os.Getenv("ASANA_ASSIGNEE"),
		project:   os.Getenv("ASANA_PROJECT"),
	}
}

func (ac *AsanaClient) CreateTask(name string) error {
	body := &addTaskRequest{
		Data: &addTaskRequestInt{
			Name:      name,
			Workspace: ac.workspace,
			Assignee:  ac.assignee,
			Projects:  []string{ac.project},
		},
	}

	js, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://app.asana.com/api/1.0/tasks", bytes.NewReader(js))
	if err != nil {
		return err
	}

	ac.addAuth(req)
	req.Header.Add("Content-Type", "application/json")

	resp, err := ac.cli.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return errors.New(resp.Status)
	}

	return nil
}

func (ac *AsanaClient) addAuth(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", ac.token))
}
