// Package keys provides API key management for HabitWire
package keys

// APIKey represents an API key
type APIKey struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Key       string `json:"key,omitempty"` // Only returned on creation
	CreatedAt string `json:"created_at,omitempty"`
	LastUsed  string `json:"last_used,omitempty"`
}

// APIKeyLean represents lean API key output
type APIKeyLean struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Key  string `json:"key,omitempty"`
}

// CreateKeyRequest represents a key creation request
type CreateKeyRequest struct {
	Name string `json:"name"`
}

// ToLean converts a full APIKey to lean output
func (k *APIKey) ToLean() APIKeyLean {
	return APIKeyLean{
		ID:   k.ID,
		Name: k.Name,
		Key:  k.Key,
	}
}

// ToLeanSlice converts a slice of APIKeys to lean output
func ToLeanSlice(keys []APIKey) []APIKeyLean {
	result := make([]APIKeyLean, len(keys))
	for i := range keys {
		result[i] = keys[i].ToLean()
	}
	return result
}
