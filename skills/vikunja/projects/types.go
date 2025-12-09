// Package projects provides project-related types and operations for Vikunja API
package projects

// Project represents a Vikunja project (based on OpenAPI spec models.Project)
type Project struct {
	ID              int64   `json:"id,omitempty"`
	Title           string  `json:"title"`
	Description     string  `json:"description,omitempty"`
	Identifier      string  `json:"identifier,omitempty"`
	HexColor        string  `json:"hex_color,omitempty"`
	ParentProjectID int64   `json:"parent_project_id,omitempty"`
	Position        float64 `json:"position,omitempty"`
	IsArchived      bool    `json:"is_archived,omitempty"`
	IsFavorite      bool    `json:"is_favorite,omitempty"`
	Created         string  `json:"created,omitempty"`
	Updated         string  `json:"updated,omitempty"`
}

// ProjectLean represents lean project output for CLI
type ProjectLean struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

// ToLean converts a full Project to lean output
func (p *Project) ToLean() ProjectLean {
	return ProjectLean{
		ID:    p.ID,
		Title: p.Title,
	}
}

// ToLeanSlice converts a slice of Projects to lean output
func ToLeanSlice(projects []Project) []ProjectLean {
	result := make([]ProjectLean, len(projects))
	for i := range projects {
		result[i] = projects[i].ToLean()
	}
	return result
}
