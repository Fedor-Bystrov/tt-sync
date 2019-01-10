package syncentity

import (
	"fmt"
	"github.com/Fedor-Bystrov/tt-sync/todoistclient"
)

// SyncEntity - grouping entity, represents
// Todoist project with all tasks and comments
type SyncEntity struct {
	Project todoistclient.Project
	Tasks   []Task
	// TODO Add version field
}

// Task - Todiost task and related comments
type Task struct {
	Task     todoistclient.Task
	Comments []todoistclient.Comment
}

func (se SyncEntity) String() string {
	return fmt.Sprintf("SyncEntity{project: %v, tasks: %v",
		se.Project, se.Tasks)
}

func (t Task) String() string {
	return fmt.Sprintf("Task{task: %v, comments: %v", t.Task, t.Comments)
}
