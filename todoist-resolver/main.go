package main

import (
	tc "github.com/Fedor-Bystrov/tt-sync/todoistclient"
	"log"
	"os"
	"regexp"
	"time"
)

var (
	todoistToken   = os.Getenv("TODOIST_TOKEN")
	goEnv          = os.Getenv("GO_ENV")
	syncPattern, _ = regexp.Compile(`^\[SYNC\]\s.+$`)
	todoistClient  = tc.NewClient(todoistToken, nil)

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
	projects := make(chan []tc.Project)
	tasks := make(chan []tc.Task)
	syncProjects := make(chan tc.Project)

	// 1. Fetching all projects and tasks
	go fetchProjects(projects)
	go fetchTasks(tasks)

	// 2. Searching for projects with [SYNC] prefix in project name
	go filterProjects(<-projects, syncProjects)

	for result := range syncProjects {
		log.Print(result)
	}
	// TODO merge projects, tasks and comments into SyncEntity
}

func filterProjects(projects []tc.Project, out chan<- tc.Project) {
	defer elapsed("[Elapsed] Searching for projects with [SYNC] prefix in project name")()
	for _, p := range projects {
		if syncPattern.MatchString(p.Name) {
			out <- p
		}
	}
	close(out)
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
	log.Print("Checking environment variables")
	if todoistToken == "" {
		log.Fatal("Fatal [Main#checkVars] Cannot resolve TODOIST_TOKEN env variable")
	}
}
