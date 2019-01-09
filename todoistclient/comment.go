package todoistclient

import "fmt"

// Comment represents Todoist Comment entity
type Comment struct {
	TaskID    uint
	ProjectID uint
	Content   string
}

func (c Comment) String() string {
	return fmt.Sprintf("Comment{task_id: %d, project_id: %d, content: %s}",
		c.TaskID, c.ProjectID, c.Content)
}
