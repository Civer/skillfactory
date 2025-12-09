// Package tasks provides task-related types and operations for Vikunja API
package tasks

// Label represents a label attached to a task
type Label struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	HexColor    string `json:"hex_color,omitempty"`
}

// Task represents a Vikunja task (based on OpenAPI spec models.Task)
type Task struct {
	ID                     int64   `json:"id,omitempty"`
	Title                  string  `json:"title"`
	Description            string  `json:"description,omitempty"`
	Done                   bool    `json:"done,omitempty"`
	DoneAt                 string  `json:"done_at,omitempty"`
	DueDate                string  `json:"due_date,omitempty"`
	StartDate              string  `json:"start_date,omitempty"`
	EndDate                string  `json:"end_date,omitempty"`
	Priority               int     `json:"priority,omitempty"`
	ProjectID              int64   `json:"project_id,omitempty"`
	RepeatAfter            int64   `json:"repeat_after,omitempty"`
	RepeatMode             int     `json:"repeat_mode,omitempty"`
	PercentDone            float64 `json:"percent_done,omitempty"`
	Identifier             string  `json:"identifier,omitempty"`
	Index                  int     `json:"index,omitempty"`
	IsFavorite             bool    `json:"is_favorite,omitempty"`
	HexColor               string  `json:"hex_color,omitempty"`
	BucketID               int64   `json:"bucket_id,omitempty"`
	Position               float64 `json:"position,omitempty"`
	CoverImageAttachmentID int64   `json:"cover_image_attachment_id,omitempty"`
	Labels                 []Label `json:"labels,omitempty"`
	Created                string  `json:"created,omitempty"`
	Updated                string  `json:"updated,omitempty"`
}

// TaskLean represents lean task output for CLI (minimal fields)
type TaskLean struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Done        bool     `json:"done"`
	Priority    int      `json:"priority"`
	DueDate     *string  `json:"due_date,omitempty"`
	StartDate   *string  `json:"start_date,omitempty"`
	EndDate     *string  `json:"end_date,omitempty"`
	ProjectID   int64    `json:"project_id"`
	Labels      []string `json:"labels,omitempty"`
	Description *string  `json:"description,omitempty"`
	PercentDone float64  `json:"percent_done,omitempty"`
	IsFavorite  bool     `json:"is_favorite,omitempty"`
}

// CreateTaskRequest represents a task creation request
type CreateTaskRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Priority    int     `json:"priority,omitempty"`
	DueDate     string  `json:"due_date,omitempty"`
	StartDate   string  `json:"start_date,omitempty"`
	EndDate     string  `json:"end_date,omitempty"`
	HexColor    string  `json:"hex_color,omitempty"`
	IsFavorite  *bool   `json:"is_favorite,omitempty"`
	PercentDone float64 `json:"percent_done,omitempty"`
}

// UpdateTaskRequest represents fields to update on a task
// Note: These are applied to an existing task before sending the full object to the API
type UpdateTaskRequest struct {
	Title       string   `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"` // Pointer to distinguish empty from unset
	Priority    *int     `json:"priority,omitempty"`
	DueDate     string   `json:"due_date,omitempty"`
	StartDate   string   `json:"start_date,omitempty"`
	EndDate     string   `json:"end_date,omitempty"`
	HexColor    string   `json:"hex_color,omitempty"`
	IsFavorite  *bool    `json:"is_favorite,omitempty"`
	PercentDone *float64 `json:"percent_done,omitempty"`
	Done        *bool    `json:"done,omitempty"`
	ProjectID   *int64   `json:"project_id,omitempty"` // For moving tasks between projects
}

// ToLean converts a full Task to lean output
func (t *Task) ToLean() TaskLean {
	// Helper to convert non-zero dates to pointers
	toDatePtr := func(date string) *string {
		if date == "" || date == "0001-01-01T00:00:00Z" {
			return nil
		}
		return &date
	}

	// Helper to convert non-empty strings to pointers
	toStrPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	// Extract label titles
	var labelTitles []string
	if len(t.Labels) > 0 {
		labelTitles = make([]string, len(t.Labels))
		for i, l := range t.Labels {
			labelTitles[i] = l.Title
		}
	}

	return TaskLean{
		ID:          t.ID,
		Title:       t.Title,
		Done:        t.Done,
		Priority:    t.Priority,
		DueDate:     toDatePtr(t.DueDate),
		StartDate:   toDatePtr(t.StartDate),
		EndDate:     toDatePtr(t.EndDate),
		ProjectID:   t.ProjectID,
		Labels:      labelTitles,
		Description: toStrPtr(t.Description),
		PercentDone: t.PercentDone,
		IsFavorite:  t.IsFavorite,
	}
}

// ToLeanSlice converts a slice of Tasks to lean output
func ToLeanSlice(tasks []Task) []TaskLean {
	result := make([]TaskLean, len(tasks))
	for i := range tasks {
		result[i] = tasks[i].ToLean()
	}
	return result
}
