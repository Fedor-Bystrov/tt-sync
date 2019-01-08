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
	testServer   *httptest.Server
	testClient   *http.Client
	wunderClient *Client

	uriMap = initURIMap()
)

func initURIMap() map[string]string {
	return map[string]string{
		// TestGetFolder response
		"/projects": `[
			{
				"id": 2200002914,
				"name": "Inbox",
				"order": 0,
				"indent": 1,
				"comment_count": 0
			},
			{
				"id": 2202928608,
				"name": "Active Projects",
				"order": 1,
				"indent": 1,
				"comment_count": 0
			}]`,
	}
}

func setUp() {
	testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Check that auth header is present
		if getHeader("X-Client-Id", req) == fmt.Sprintf("Bearer %s", authToken) {
			res.Write([]byte(uriMap[req.RequestURI]))
		} else {
			res.Write(nil)
		}
	}))
	testClient = testServer.Client()
	wunderClient = NewClient(authToken, testClient)
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
	assert.NotNil(t, nil)
}
