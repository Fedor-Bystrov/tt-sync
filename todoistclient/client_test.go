package todoistclient

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
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
				"url": "url_1"
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
				"url": "url_2",
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
				"url": "url_3"
			}]`,
		// TestGetComments response
		"/comments?task_id=123": `[
				{
					"id": 1,
					"task_id": 11,
					"posted": "2019-01-09T10:50:59Z",
					"content": "F"
				},
				{
					"id": 2,
					"task_id": 21,
					"posted": "2019-01-09T10:51:02Z",
					"content": "g"
				},
				{
					"id": 3,
					"task_id": 31,
					"posted": "2019-01-09T10:51:02Z",
					"content": "C"
				}]`,
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
		Task{ID: 3, ProjectID: 1, Content: "content_1", Completed: false,
			Order: 3, Indent: 1, Priority: 1, CommentCount: 0, URL: "url_1", LabelIDs: make([]uint, 0)},

		Task{ID: 4, ProjectID: 2, Content: "content_2", Completed: false, Due: &Due{false, "Jan 8", "2019-01-08"},
			Order: 4, Indent: 1, Priority: 1, CommentCount: 0, URL: "url_2", LabelIDs: []uint{1}},

		Task{ID: 5, ProjectID: 3, Content: "content_3", Completed: true,
			Order: 4, Indent: 1, Priority: 1, CommentCount: 10, URL: "url_3", LabelIDs: []uint{1, 2, 3}},
	}
	tasks, err := todoistClient.GetTasks()
	assert.Nil(t, err)
	assert.ElementsMatch(t, expTasks, tasks)
}

func TestGetComments(t *testing.T) {
	posted1, _ := time.Parse(time.RFC3339, "2019-01-09T10:50:59Z")
	posted2, _ := time.Parse(time.RFC3339, "2019-01-09T10:51:02Z")
	posted3, _ := time.Parse(time.RFC3339, "2019-01-09T10:51:02Z")
	expComments := []Comment{
		Comment{ID: 1, TaskID: 11, Content: "F", Posted: posted1},
		Comment{ID: 2, TaskID: 21, Content: "g", Posted: posted2},
		Comment{ID: 3, TaskID: 31, Content: "C", Posted: posted3},
	}
	comments, err := todoistClient.GetComments(123)
	fmt.Println("Tasks: ", comments)
	assert.Nil(t, err)
	assert.ElementsMatch(t, expComments, comments)
}
