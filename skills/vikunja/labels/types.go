// Package labels provides label-related types and operations for Vikunja API
package labels

// Label represents a Vikunja label (based on OpenAPI spec models.Label)
type Label struct {
	ID          int64  `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	HexColor    string `json:"hex_color,omitempty"`
	Created     string `json:"created,omitempty"`
	Updated     string `json:"updated,omitempty"`
}

// LabelLean represents lean label output for CLI
type LabelLean struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	HexColor string `json:"hex_color"`
}

// CreateLabelRequest represents a label creation request
type CreateLabelRequest struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	HexColor    string `json:"hex_color,omitempty"`
}

// UpdateLabelRequest represents a label update request
type UpdateLabelRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	HexColor    string `json:"hex_color,omitempty"`
}

// ToLean converts a full Label to lean output
func (l *Label) ToLean() LabelLean {
	return LabelLean{
		ID:       l.ID,
		Title:    l.Title,
		HexColor: l.HexColor,
	}
}

// ToLeanSlice converts a slice of Labels to lean output
func ToLeanSlice(labels []Label) []LabelLean {
	result := make([]LabelLean, len(labels))
	for i := range labels {
		result[i] = labels[i].ToLean()
	}
	return result
}
