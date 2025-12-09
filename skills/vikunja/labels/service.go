// Package labels provides label operations for Vikunja API
package labels

import (
	"encoding/json"
	"fmt"

	"github.com/petervogelmann/skillfactory/skills/vikunja/client"
)

// Service handles label operations
type Service struct {
	client *client.Client
}

// NewService creates a new label service
func NewService(c *client.Client) *Service {
	return &Service{client: c}
}

// List retrieves all labels
func (s *Service) List() ([]Label, error) {
	data, err := s.client.Get("/labels")
	if err != nil {
		return nil, err
	}

	var labels []Label
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, fmt.Errorf("failed to parse labels: %w", err)
	}

	return labels, nil
}

// Get retrieves a single label by ID
func (s *Service) Get(labelID int64) (*Label, error) {
	endpoint := fmt.Sprintf("/labels/%d", labelID)

	data, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("failed to parse label: %w", err)
	}

	return &label, nil
}

// Create creates a new label
func (s *Service) Create(req CreateLabelRequest) (*Label, error) {
	data, err := s.client.Put("/labels", req)
	if err != nil {
		return nil, err
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("failed to parse created label: %w", err)
	}

	return &label, nil
}

// Update updates an existing label
func (s *Service) Update(labelID int64, req UpdateLabelRequest) (*Label, error) {
	endpoint := fmt.Sprintf("/labels/%d", labelID)

	data, err := s.client.Post(endpoint, req)
	if err != nil {
		return nil, err
	}

	var label Label
	if err := json.Unmarshal(data, &label); err != nil {
		return nil, fmt.Errorf("failed to parse updated label: %w", err)
	}

	return &label, nil
}

// Delete deletes a label
func (s *Service) Delete(labelID int64) error {
	endpoint := fmt.Sprintf("/labels/%d", labelID)
	_, err := s.client.Delete(endpoint)
	return err
}
