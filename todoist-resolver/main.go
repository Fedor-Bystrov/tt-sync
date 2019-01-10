package main

import (
	"github.com/Fedor-Bystrov/tt-sync/todoistclient"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	todoistToken   = os.Getenv("TODOIST_TOKEN")
	goEnv          = os.Getenv("GO_ENV")
	syncPattern, _ = regexp.Compile(`^\[SYNC\]\s.+$`)

	logFile *os.File
)

func init() {
	defer elapsed("[Elapsed] INIT")()
	if goEnv != "development" {
		logFile, err := os.OpenFile("todoist-resolver.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Fatal [Main#init] Error opening file: %v", err)
		}
		log.SetOutput(logFile)
	}
}

func main() {
	defer elapsed("[Elapsed] MAIN")()
	if goEnv != "development" {
		defer logFile.Close()
	}
	checkVars(todoistToken)

	todoistClient := todoistclient.NewClient(todoistToken, nil)
	projects, err := todoistClient.GetProjects()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch projects", err)
	}

	var syncProject todoistclient.Project
	for _, p := range projects {
		if syncPattern.MatchString(p.Name) {
			syncProject = p
		}
	}

	tasks, err := todoistClient.GetTasks()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch tasks", err)
	}

	for _, t := range tasks {
		if t.ProjectID == syncProject.ID {
			log.Print(t)
		}
	}
}

func elapsed(message string) func() {
	start := time.Now()
	return func() { log.Printf("%s took %v\n", message, time.Since(start)) }
}

func checkVars(todoistToken string) {
	log.Print("Checking environment variables")
	if todoistToken == "" {
		log.Fatal("Fatal [Main#checkVars] Cannot resolve TODOIST_TOKEN env variable")
	}
}
