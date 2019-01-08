package todoistclient

import (
	"fmt"
	"net/http"
	"time"
)

var (
	apiURL = "https://beta.todoist.com/API/v8"

	httpClient *http.Client
)

// Client is a Todoist rest api client
type Client struct {
	token string
}

// NewClient returns a new instance of todoist client
// given todoist app token and, optionally, http.Client
func NewClient(token string, client *http.Client) *Client {
	if client == nil {
		httpClient = &http.Client{Timeout: time.Second * 10}
	} else {
		httpClient = client
	}
	return &Client{token}
}

// GetProjects returns all projects for user
// corresponding to given token
func (c Client) GetProjects() []Project {
	return nil
}

// GetTasks returns all tasks for user
// corresponding to given token
func (c Client) GetTasks() []Task {
	return nil
}

// GetComments returns all comments for user
// corresponding to given token
func (c Client) GetComments() []Comment {
	return nil
}

func (c Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return req, nil
}
