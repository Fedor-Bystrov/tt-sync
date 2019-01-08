package todoistclient

import "fmt"

// Project represents Todoist Project entity
// https://developer.todoist.com/rest/v8/#projects
type Project struct {
	ID           uint
	Name         string
	Order        uint
	Indent       uint
	CommentCount uint `json:"comment_count"`
}

func (p Project) String() string {
	return fmt.Sprintf("Project{id: %d, name: %s, order: %d, indent: %d, comment_count: %d}",
		p.ID, p.Name, p.Order, p.Indent, p.CommentCount)
}
