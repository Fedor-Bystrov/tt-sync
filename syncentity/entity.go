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

// AddTask add task to SyncEntity tasks slice
// if slice is undefined, creates new slice
// and appends task to it
func (se *SyncEntity) AddTask(t Task) {
	if se.Tasks == nil {
		se.Tasks = make([]Task, 0)
	}
	se.Tasks = append(se.Tasks, t)
}

func (se SyncEntity) String() string {
	return fmt.Sprintf("SyncEntity{project: %v, tasks: %v",
		se.Project, se.Tasks)
}

// Task - Todiost task and related comments
type Task struct {
	Task     todoistclient.Task
	Comments []todoistclient.Comment
}

func (t Task) String() string {
	return fmt.Sprintf("Task{task: %v, comments: %v", t.Task, t.Comments)
}
