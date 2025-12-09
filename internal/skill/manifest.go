// Package skill handles skill manifest parsing and discovery
package skill

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Variable represents a configurable skill variable
type Variable struct {
	Name        string `yaml:"name"`
	Label       string `yaml:"label"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Placeholder string `yaml:"placeholder"`
	Default     string `yaml:"default"`
	Type        string `yaml:"type"` // string, secret, json
}

// BuildConfig holds build configuration
type BuildConfig struct {
	Entry  string `yaml:"entry"`
	Binary string `yaml:"binary"`
}

// DeployFile represents a file to deploy
type DeployFile struct {
	Source string `yaml:"source"`
	Target string `yaml:"target"`
}

// DeployConfig holds deploy configuration
type DeployConfig struct {
	Files   []DeployFile `yaml:"files"`
	Wrapper bool         `yaml:"wrapper"`
}

// DocsConfig holds documentation configuration
type DocsConfig struct {
	Template string `yaml:"template"`
	Output   string `yaml:"output"`
}

// Manifest represents a skill.yaml file
type Manifest struct {
	Name             string       `yaml:"name"`
	Description      string       `yaml:"description"`
	SkillDescription string       `yaml:"skill_description"` // Optional: longer description for SKILL.md frontmatter
	Version          string       `yaml:"version"`
	Variables        []Variable   `yaml:"variables"`
	Build            BuildConfig  `yaml:"build"`
	Deploy           DeployConfig `yaml:"deploy"`
	Docs             DocsConfig   `yaml:"docs"`

	// Runtime fields (not from YAML)
	Path string `yaml:"-"` // Path to skill directory
}

// GetSkillDescription returns SkillDescription if set, otherwise Description
func (m *Manifest) GetSkillDescription() string {
	if m.SkillDescription != "" {
		return m.SkillDescription
	}
	return m.Description
}

// SkillError represents a skill that failed to load
type SkillError struct {
	Name  string
	Path  string
	Error error
}

// LoadManifest loads a skill manifest from a directory
func LoadManifest(skillDir string) (*Manifest, error) {
	manifestPath := filepath.Join(skillDir, "skill.yaml")

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill.yaml: %w", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse skill.yaml: %w", err)
	}

	manifest.Path = skillDir
	return &manifest, nil
}

// DiscoverSkills finds all skills in a directory
// Returns valid manifests and a list of skills that failed to load
func DiscoverSkills(baseDir string) ([]*Manifest, []SkillError, error) {
	skillsDir := filepath.Join(baseDir, "skills")

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read skills directory: %w", err)
	}

	var manifests []*Manifest
	var errors []SkillError

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillDir := filepath.Join(skillsDir, entry.Name())
		manifestPath := filepath.Join(skillDir, "skill.yaml")

		// Check if skill.yaml exists
		if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
			continue
		}

		manifest, err := LoadManifest(skillDir)
		if err != nil {
			errors = append(errors, SkillError{
				Name:  entry.Name(),
				Path:  skillDir,
				Error: err,
			})
			continue
		}

		manifests = append(manifests, manifest)
	}

	return manifests, errors, nil
}

// GetRequiredVariables returns only required variables
func (m *Manifest) GetRequiredVariables() []Variable {
	var required []Variable
	for _, v := range m.Variables {
		if v.Required {
			required = append(required, v)
		}
	}
	return required
}

// GetOptionalVariables returns only optional variables
func (m *Manifest) GetOptionalVariables() []Variable {
	var optional []Variable
	for _, v := range m.Variables {
		if !v.Required {
			optional = append(optional, v)
		}
	}
	return optional
}
