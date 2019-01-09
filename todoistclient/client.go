package todoistclient

import (
	"encoding/json"
	"fmt"
	"log"
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
func (c Client) GetProjects() ([]Project, error) {
	log.Print("[TodoistClient#GetProjects] Fetching projects")
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/projects", apiURL))
	if err != nil {
		log.Print("[TodoistClient#GetProjects] Error fetching projects", err)
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Print("[TodoistClient#GetProjects] Error fetching projects", err)
		return nil, err
	}
	defer res.Body.Close()

	projects := make([]Project, 0)
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		log.Print("[TodoistClient#GetProjects] Error during decoding todoist response", err)
		return nil, err
	}

	log.Print("[TodoistClient#GetProjects] Projects fetched successfully")
	return projects, nil
}

// GetTasks returns all tasks for user
// corresponding to given token
func (c Client) GetTasks() ([]Task, error) {
	log.Print("[TodoistClient#GetTasks] Fetching tasks")
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/tasks", apiURL))
	if err != nil {
		log.Print("[TodoistClient#GetTasks] Error fetching tasks", err)
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Print("[TodoistClient#GetTasks] Error fetching tasks", err)
		return nil, err
	}
	defer res.Body.Close()

	tasks := make([]Task, 0)
	err = json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		log.Print("[TodoistClient#GetTasks] Error during decoding todoist response", err)
		return nil, err
	}

	log.Print("[TodoistClient#GetTasks] Tasks fetched successfully")
	return tasks, nil
}

// GetComments returns all task comments given task_id
func (c Client) GetComments(taskID uint) ([]Comment, error) {
	log.Print(fmt.Sprintf("[TodoistClient#GetComments] Fetching comments for task_id: %d", taskID))
	req, err := c.newRequest(http.MethodGet, fmt.Sprintf("%s/comments?task_id=%d", apiURL, taskID))
	if err != nil {
		log.Print("[TodoistClient#GetComments] Error fetching comments", err)
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		log.Print("[TodoistClient#GetComments] Error fetching comments", err)
		return nil, err
	}
	defer res.Body.Close()

	comments := make([]Comment, 0)
	err = json.NewDecoder(res.Body).Decode(&comments)
	if err != nil {
		log.Print("[TodoistClient#GetComments] Error during decoding todoist response", err)
		return nil, err
	}

	log.Print("[TodoistClient#GetComments] Comments fetched successfully")
	return comments, nil
}

func (c Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	return req, nil
}
