// Package keys provides API key operations for HabitWire API
package keys

import (
	"encoding/json"
	"fmt"

	"habitwire/client"
)

// Service handles API key operations
type Service struct {
	client *client.Client
}

// NewService creates a new API key service
func NewService(c *client.Client) *Service {
	return &Service{client: c}
}

// List retrieves all API keys
func (s *Service) List() ([]APIKey, error) {
	data, err := s.client.Get("/keys")
	if err != nil {
		return nil, err
	}

	var keys []APIKey
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse keys: %w", err)
	}

	return keys, nil
}

// Create creates a new API key
func (s *Service) Create(req CreateKeyRequest) (*APIKey, error) {
	data, err := s.client.Post("/keys", req)
	if err != nil {
		return nil, err
	}

	var key APIKey
	if err := json.Unmarshal(data, &key); err != nil {
		return nil, fmt.Errorf("failed to parse created key: %w", err)
	}

	return &key, nil
}

// Delete deletes an API key
func (s *Service) Delete(keyID string) error {
	endpoint := fmt.Sprintf("/keys/%s", keyID)
	_, err := s.client.Delete(endpoint)
	return err
}
