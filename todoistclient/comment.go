package todoistclient

import (
	"fmt"
	"time"
)

// Comment represents Todoist Comment entity
type Comment struct {
	ID        uint
	TaskID    uint `json:"task_id"`
	ProjectID uint `json:"project_id"`
	Posted    time.Time
	Content   string
}

func (c Comment) String() string {
	return fmt.Sprintf("Comment{id: %d, task_id: %d, project_id: %d, content: %s, posted: %v}",
		c.ID, c.TaskID, c.ProjectID, c.Content, c.Posted)
}
