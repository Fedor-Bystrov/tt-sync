package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	se "github.com/Fedor-Bystrov/tt-sync/syncentity"
	tc "github.com/Fedor-Bystrov/tt-sync/todoistclient"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/readpref"
)

var (
	todoistToken   = os.Getenv("TODOIST_TOKEN")
	goEnv          = os.Getenv("GO_ENV")
	syncPattern, _ = regexp.Compile(`^\[SYNC\]\s.+$`)
	todoistClient  = tc.NewClient(todoistToken, nil)

	ctx, cancelFunc = context.WithTimeout(context.Background(), 10*time.Second)
	client          mongo.Client
	logFile         *os.File
)

func init() {
	defer elapsed("[Elapsed] INIT")()
	checkVars(todoistToken)

	log.Print("[Main#init] Connecting to mongo")
	client, err := mongo.Connect(ctx, "mongodb://localhost:27017")

	// Set up logging
	if goEnv != "development" {
		logFile, err := os.OpenFile("todoist-resolver.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Fatal [Main#init] Error opening file: %v", err)
		}
		log.SetOutput(logFile)
	}

	// Check that connection established
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("Fatal [Main#init] Error connecting to mongo: %v", err)
	}
	log.Print("[Main#init] Connection established")
}

func main() {
	defer elapsed("[Elapsed] MAIN")()
	if goEnv != "development" {
		defer logFile.Close()
	}

	projectRespCh := make(chan []tc.Project)
	tasksRespCh := make(chan []tc.Task)
	tasksCh := make(chan se.Task)

	// 1. Fetching all projects and tasks
	go fetchProjects(projectRespCh)
	go fetchTasks(tasksRespCh)

	// 2. Searching for projects with [SYNC] prefix in project name
	projects := make([]tc.Project, 0)
	for _, p := range <-projectRespCh {
		if syncPattern.MatchString(p.Name) {
			projects = append(projects, p)
		}
	}

	// 3. Fetching comments for every task in each sync project
	tasks := make([]se.Task, 0)
	go resolveSyncTasks(<-tasksRespCh, projects, tasksCh)
	for t := range tasksCh {
		tasks = append(tasks, t)
	}

	// 4. Merging projects and tasks into sync_entities
	entities := make([]se.SyncEntity, len(projects))
	for _, p := range projects {
		e := se.SyncEntity{Project: p}
		for _, t := range tasks {
			if p.ID == t.Task.ProjectID {
				e.AddTask(t)
			}
		}
		entities = append(entities, e)
	}

	for _, e := range entities {
		log.Print(e)
	}
}

func resolveSyncTasks(tasks []tc.Task, ps []tc.Project, out chan se.Task) {
	defer elapsed("[Elapsed] MAIN#resolveSyncTasks goroutine")()
	var wg sync.WaitGroup
	for _, p := range ps {
		for _, t := range tasks {
			if t.ProjectID == p.ID {
				wg.Add(1)
				go fetchTaskComments(t, out, &wg)
			}
		}
	}
	wg.Wait()
	close(out)
}

func fetchTaskComments(task tc.Task, out chan se.Task, wg *sync.WaitGroup) {
	defer elapsed(fmt.Sprintf("[Elapsed] MAIN#fetchTaskComments goroutine for task_id: %d", task.ID))()
	comments, err := todoistClient.GetComments(task.ID)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fatal [Main] cannot fetch comments for task_id: %d, %v", task.ID, err))
	}
	out <- se.Task{Task: task, Comments: comments}
	wg.Done()
}

func fetchProjects(out chan<- []tc.Project) {
	defer elapsed("[Elapsed] MAIN#fetchProjects goroutine")()
	projects, err := todoistClient.GetProjects()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch projects", err)
	}
	out <- projects
	close(out)
}

func fetchTasks(out chan<- []tc.Task) {
	defer elapsed("[Elapsed] MAIN#fetchTasks goroutine")()
	tasks, err := todoistClient.GetTasks()
	if err != nil {
		log.Fatal("Fatal [Main] cannot fetch tasks", err)
	}
	out <- tasks
	close(out)
}

func elapsed(message string) func() {
	start := time.Now()
	return func() { log.Printf("%s took %v\n", message, time.Since(start)) }
}

func checkVars(todoistToken string) {
	log.Print("[Main#init] Checking environment variables")
	if todoistToken == "" {
		log.Fatal("Fatal [Main#checkVars] Cannot resolve TODOIST_TOKEN env variable")
	}
}
