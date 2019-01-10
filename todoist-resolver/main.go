package main

import (
	"github.com/Fedor-Bystrov/tt-sync/todoistclient"
	"log"
	"os"
	"regexp"
)

var (
	todoistToken   = os.Getenv("TODOIST_TOKEN")
	goEnv          = os.Getenv("GO_ENV")
	syncPattern, _ = regexp.Compile(`^\[SYNC\]\s.+$`)

	logFile *os.File
)

func init() {
	if goEnv != "development" {
		logFile, err := os.OpenFile("todoist-resolver.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("Fatal [Main#init] Error opening file: %v", err)
		}
		log.SetOutput(logFile)
	}
}

func main() {
	if goEnv != "development" {
		defer logFile.Close()
	}
	checkVars(todoistToken)

	todoistClient := todoistclient.NewClient(todoistToken, nil)
	projects, err := todoistClient.GetProjects()
	if err != nil {
		log.Fatal("Fatal [Main#GetProjects] cannot fetch projects", err)
	}

	for _, p := range projects {
		if syncPattern.MatchString(p.Name) {
			log.Print(p)
		}
	}
}

func checkVars(todoistToken string) {
	log.Print("Checking environment variables")
	if todoistToken == "" {
		log.Fatal("Fatal [Main#checkVars] Cannot resolve TODOIST_TOKEN env variable")
	}
}
