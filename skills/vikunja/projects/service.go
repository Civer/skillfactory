// Package projects provides project operations for Vikunja API
package projects

import (
	"encoding/json"
	"fmt"

	"github.com/petervogelmann/skillfactory/skills/vikunja/client"
)

// Service handles project operations
type Service struct {
	client *client.Client
}

// NewService creates a new project service
func NewService(c *client.Client) *Service {
	return &Service{client: c}
}

// List retrieves all projects
func (s *Service) List() ([]Project, error) {
	data, err := s.client.Get("/projects")
	if err != nil {
		return nil, err
	}

	var projects []Project
	if err := json.Unmarshal(data, &projects); err != nil {
		return nil, fmt.Errorf("failed to parse projects: %w", err)
	}

	return projects, nil
}

// Get retrieves a single project by ID
func (s *Service) Get(projectID int64) (*Project, error) {
	endpoint := fmt.Sprintf("/projects/%d", projectID)

	data, err := s.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, fmt.Errorf("failed to parse project: %w", err)
	}

	return &project, nil
}
