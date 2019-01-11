package main

import (
	"github.com/Fedor-Bystrov/tt-sync/syncentity"
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
	todoistClient  = todoistclient.NewClient(todoistToken, nil)

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
	projectsChan := make(chan []todoistclient.Project)
	tasksChan := make(chan []todoistclient.Task)

	// 1. Fetching all projects and tasks
	go fetchProjects(projectsChan)
	go fetchTasks(tasksChan)

	// 2. Searching for projects with [SYNC] prefix in project name
	syncEntities := make([]syncentity.SyncEntity, 0)
	for _, p := range <-projectsChan {
		if syncPattern.MatchString(p.Name) {
			syncEntities = append(syncEntities, syncentity.SyncEntity{
				Project: p,
				Tasks:   make([]syncentity.Task, 0),
			})
		}
	}

	// 4. Merging sync project, related tasks and comments into one entity
	for i := range syncEntities {
		for _, t := range <-tasksChan {
			if syncEntities[i].Project.ID == t.ProjectID {
				comments, err := todoistClient.GetComments(t.ID)
				if err != nil {
					log.Fatalf("Fatal [Main] cannot fetch comments for task_id: %d, err: %v", t.ID, err)
				}
				syncEntities[i].Tasks = append(syncEntities[i].Tasks, syncentity.Task{
					Task:     t,
					Comments: comments,
				})
			}
		}
	}

	for _, s := range syncEntities {
		log.Print(s)
	}
}

func fetchProjects(projectsChan chan<- []todoistclient.Project) {
	defer elapsed("[Elapsed] MAIN#fetchProjects goroutine")()
	projects, err := todoistClient.GetProjects()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch projects", err)
	}
	projectsChan <- projects
	close(projectsChan)
}

func fetchTasks(tasksChan chan<- []todoistclient.Task) {
	defer elapsed("[Elapsed] MAIN#fetchTasks goroutine")()
	tasks, err := todoistClient.GetTasks()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch tasks", err)
	}
	tasksChan <- tasks
	close(tasksChan)
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
