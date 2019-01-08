package todoistclient

import (
	"fmt"
	"time"
)

// Task is a Todoist Task entity
// Properties:
// ID - Task id.
// ProjectID - Integer Task’s project id (read-only).
// Content  - String Task content.
// Completed - Flag to mark completed tasks.
// LabelIDs - Array of label ids, associated with a task.
// Order - Position in the project (read-only).
// Indent - Task indentation level from 1 to 5 (read-only).
// Priority - Task priority from 1 (normal, default value) to 4 (urgent).
// Due - object representing task due date/time (described below).
// URL - URL to access this task in Todoist web interface.
// CommentCount - Number of task comments.
type Task struct {
	ID           uint
	ProjectID    uint `json:"project_id"`
	Content      string
	Completed    bool
	LabelIDs     []uint `json:"label_ids"`
	Order        uint
	Indent       uint
	Priority     uint
	Due          time.Time
	URL          string
	CommentCount uint `json:"comment_count"`
}

func (t Task) String() string {
	return fmt.Sprintf("Task{id: %d project_id: %d, content: %s, url: %s}",
		t.ID, t.ProjectID, t.Content, t.URL)
}
