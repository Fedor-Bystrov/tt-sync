package todoistclient

import (
	"fmt"
)

// Task represents Todoist Task entity
// Properties:
// ID - Task id.
// ProjectID - Integer Task’s project id (read-only).
// Content  - String Task content.
// Completed - Flag to mark completed tasks.
// LabelIDs - Array of label ids, associated with a task.
// Order - Position in the project (read-only).
// Indent - Task indentation level from 1 to 5 (read-only).
// Priority - Task priority from 1 (normal, default value) to 4 (urgent).
// Due - object representing task due date/time.
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
	Due          *Due
	URL          string
	CommentCount uint `json:"comment_count"`
}

func (t Task) String() string {
	return fmt.Sprintf("Task{id: %d, project_id: %d, content: %s, completed: %v, label_ids: %v, "+
		"order: %d indent: %d, priority: %d, due: %v, url: %s, comment_count: %d}",
		t.ID, t.ProjectID, t.Content, t.Completed, t.LabelIDs, t.Order, t.Indent,
		t.Priority, t.Due, t.URL, t.CommentCount)
}

// Due represents Todoist Due date object
type Due struct {
	Recurring bool
	Str       string `json:"string"`
	Date      string
}

func (d Due) String() string {
	return fmt.Sprintf("Due{reccuring: %v string: %s, date: %s}", d.Recurring, d.Str, d.Date)
}
