package todoistclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

const (
	authToken = "TodoistToken"
)

var (
	testServer    *httptest.Server
	testClient    *http.Client
	todoistClient *Client

	uriMap = initURIMap()
)

func initURIMap() map[string]string {
	return map[string]string{
		// TestGetProjects response
		"/projects": `[
			{
				"id": 1,
				"name": "Inbox",
				"order": 0,
				"indent": 1,
				"comment_count": 0
			},
			{
				"id": 2,
				"name": "Active Projects",
				"order": 1,
				"indent": 1,
				"comment_count": 0
			}]`,
		// TestGetTasks response
		"/tasks": `[
			{
				"id": 3,
				"project_id": 1,
				"content": "content_1",
				"completed": false,
				"label_ids": [],
				"order": 3,
				"indent": 1,
				"priority": 1,
				"comment_count": 0,
				"url": "https://todoist.com/showTask?id=3"
			},
			{
				"id": 4,
				"project_id": 2,
				"content": "content_2",
				"completed": false,
				"label_ids": [1],
				"order": 4,
				"indent": 1,
				"priority": 1,
				"comment_count": 0,
				"url": "https://todoist.com/showTask?id=4",
				"due": {
					"recurring": false,
					"string": "Jan 8",
					"date": "2019-01-08"
				}
			},
			{
				"id": 5,
				"project_id": 3,
				"content": "content_3",
				"completed": true,
				"label_ids": [1,2,3],
				"order": 4,
				"indent": 1,
				"priority": 1,
				"comment_count": 10,
				"url": "https://todoist.com/showTask?id=5"
			}`,
	}
}

func setUp() {
	testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check that auth header is present
		if getHeader("Authorization", req) == fmt.Sprintf("Bearer %s", authToken) {
			res.Write([]byte(uriMap[req.RequestURI]))
		} else {
			res.Write(nil)
		}
	}))
	testClient = testServer.Client()
	todoistClient = NewClient(authToken, testClient)
	// Changing apiURL to point to testServer
	reflect.ValueOf(&apiURL).Elem().SetString(testServer.URL)
}

func getHeader(key string, req *http.Request) string {
	value := req.Header[key]
	if len(value) > 0 {
		return value[0]
	}
	return ""
}

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	testServer.Close()
	os.Exit(retCode)
}

func TestGetProjects(t *testing.T) {
	expProjects := []Project{
		Project{1, "Inbox", 0, 1, 0},
		Project{2, "Active Projects", 1, 1, 0},
	}
	projects, err := todoistClient.GetProjects()
	assert.Nil(t, err)
	assert.ElementsMatch(t, expProjects, projects)
}

func TestGetTasks(t *testing.T) {
	expTasks := []Task{
		Task{3, 1, "content_1", false, nil, 3, 1, 1, nil,
			"\"https://todoist.com/showTask?id=3\"", 0},
		Task{4, 2, "content_2", false, []uint{1}, 4, 1, 1,
			&Due{false, "Jan 8", "2019-01-08"},
			"\"https://todoist.com/showTask?id=4\"", 0},
		Task{5, 3, "content_3", true, []uint{1, 2, 3}, 4, 1, 1, nil,
			"\"https://todoist.com/showTask?id=5\"", 10},
	}
	tasks, err := todoistClient.GetTasks()
	assert.Nil(t, err)
	assert.ElementsMatch(t, expTasks, tasks)
}
