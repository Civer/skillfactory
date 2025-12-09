// Package tui provides the terminal user interface for SkillFactory
package tui

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/petervogelmann/skillfactory/internal/skill"
)

// View represents different screens in the TUI
type View int

const (
	ViewSkillList View = iota // Start here: list of available skills
	ViewConfig                // Configure selected skill
	ViewConfirm               // Confirm and deploy
	ViewOverwrite             // Warning: skill already exists
	ViewBuilding              // Building in progress
	ViewDone                  // Success/Error result
)

// Model represents the application state
type Model struct {
	projectRoot   string
	version       string
	manifests     []*skill.Manifest
	skillErrors   []skill.SkillError // Skills that failed to load
	currentView   View
	skillCursor   int
	selectedSkill *skill.Manifest
	selectedError *skill.SkillError // Selected error skill (for viewing errors)

	// Dynamic input fields based on skill variables
	inputs       []textinput.Model
	inputLabels  []string
	inputFocus   int

	// Deploy configuration
	skillsFolder    string // Base folder for skills (e.g., /path/to/.claude/skills/)
	skillFolderName string // Subfolder name for this skill (default: skill name)

	// Configured values
	configValues map[string]string

	// Status messages
	statusMsg string
	errorMsg  string

	// Build state
	building    bool
	buildOutput string

	width    int
	height   int
	quitting bool
}

// NewModel creates a new TUI model
func NewModel(projectRoot string, version string) Model {
	// Discover skills
	manifests, skillErrors, err := skill.DiscoverSkills(projectRoot)
	if err != nil {
		// Will show error in UI
		manifests = []*skill.Manifest{}
		skillErrors = []skill.SkillError{}
	}

	return Model{
		projectRoot:  projectRoot,
		version:      version,
		manifests:    manifests,
		skillErrors:  skillErrors,
		currentView:  ViewSkillList,
		configValues: make(map[string]string),
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		// Config view: handle inputs FIRST, then navigation
		if m.currentView == ViewConfig {
			switch msg.String() {
			case "esc":
				m.currentView = ViewSkillList
				m.errorMsg = ""
				return m, nil
			case "tab", "down":
				m.inputs[m.inputFocus].Blur()
				m.inputFocus = (m.inputFocus + 1) % len(m.inputs)
				m.inputs[m.inputFocus].Focus()
				return m, textinput.Blink
			case "shift+tab", "up":
				m.inputs[m.inputFocus].Blur()
				m.inputFocus--
				if m.inputFocus < 0 {
					m.inputFocus = len(m.inputs) - 1
				}
				m.inputs[m.inputFocus].Focus()
				return m, textinput.Blink
			case "ctrl+d", "enter":
				// Done editing - validate and continue
				if m.validateInputs() {
					m.saveConfigValues()
					m.currentView = ViewConfirm
				}
				return m, nil
			default:
				// Pass ALL other keys to the text input
				return m.updateInputs(msg)
			}
		}

		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case buildCompleteMsg:
		m.building = false
		m.buildOutput = msg.output
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
			m.currentView = ViewDone
		} else {
			// Build succeeded, now deploy
			return m, m.deploySkill()
		}
		return m, nil

	case deployCompleteMsg:
		m.currentView = ViewDone
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.statusMsg = "Skill deployed successfully!"
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.currentView {
	case ViewSkillList:
		return m.handleSkillListView(msg)
	case ViewConfirm:
		return m.handleConfirmView(msg)
	case ViewOverwrite:
		return m.handleOverwriteView(msg)
	case ViewDone:
		return m.handleDoneView(msg)
	}
	return m, nil
}

func (m Model) handleSkillListView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	totalItems := len(m.manifests) + len(m.skillErrors)

	switch msg.String() {
	case "q", "esc":
		m.quitting = true
		return m, tea.Quit
	case "up", "k":
		if m.skillCursor > 0 {
			m.skillCursor--
		}
	case "down", "j":
		if m.skillCursor < totalItems-1 {
			m.skillCursor++
		}
	case "enter":
		if m.skillCursor < len(m.manifests) {
			// Valid skill selected
			m.selectedSkill = m.manifests[m.skillCursor]
			m.selectedError = nil
			m.currentView = ViewConfig
			m.setupInputsFromManifest()
			return m, textinput.Blink
		} else if m.skillCursor < totalItems {
			// Error skill selected - show error details
			errorIdx := m.skillCursor - len(m.manifests)
			m.selectedError = &m.skillErrors[errorIdx]
			m.selectedSkill = nil
			m.errorMsg = m.selectedError.Error.Error()
		}
	}
	return m, nil
}

func (m Model) handleConfirmView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		m.currentView = ViewConfig
		m.inputFocus = 0
		if len(m.inputs) > 0 {
			m.inputs[0].Focus()
		}
		return m, textinput.Blink
	case "enter", "y":
		// Check if skill already exists
		if m.skillExists() {
			m.currentView = ViewOverwrite
			return m, nil
		}
		// Start build
		m.currentView = ViewBuilding
		m.building = true
		m.errorMsg = ""
		m.statusMsg = ""
		return m, m.startBuild()
	}
	return m, nil
}

func (m Model) handleOverwriteView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "n":
		// Go back to confirm view
		m.currentView = ViewConfirm
		return m, nil
	case "y":
		// Proceed with build (overwrite)
		m.currentView = ViewBuilding
		m.building = true
		m.errorMsg = ""
		m.statusMsg = ""
		return m, m.startBuild()
	}
	return m, nil
}

// skillExists checks if the deploy path already contains a skill
func (m Model) skillExists() bool {
	deployPath := m.getDeployPath()
	if deployPath == "" {
		return false
	}

	// Check if the bin directory with binary exists
	binaryName := m.selectedSkill.Build.Binary
	if binaryName == "" {
		binaryName = m.selectedSkill.Name
	}
	binaryPath := filepath.Join(deployPath, "bin", binaryName)

	_, err := os.Stat(binaryPath)
	return err == nil
}

func (m Model) handleDoneView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "q", "esc":
		m.quitting = true
		return m, tea.Quit
	case "r":
		// Restart - go back to skill list
		m.currentView = ViewSkillList
		m.errorMsg = ""
		m.statusMsg = ""
		m.buildOutput = ""
		return m, nil
	}
	return m, nil
}

// setupInputsFromManifest creates input fields from the skill's variables
func (m *Model) setupInputsFromManifest() {
	if m.selectedSkill == nil {
		return
	}

	// Create inputs for all variables + 2 deploy fields (skills folder + skill name)
	numInputs := len(m.selectedSkill.Variables) + 2
	m.inputs = make([]textinput.Model, numInputs)
	m.inputLabels = make([]string, numInputs)

	for i, v := range m.selectedSkill.Variables {
		input := textinput.New()
		input.Placeholder = v.Placeholder
		if input.Placeholder == "" && v.Default != "" {
			input.Placeholder = v.Default
		}
		input.CharLimit = 200
		input.Width = 50

		if v.Type == "secret" {
			input.EchoMode = textinput.EchoPassword
		}

		// Load existing value if any (but don't pre-fill defaults)
		if val, ok := m.configValues[v.Name]; ok {
			input.SetValue(val)
		}
		// Default is only shown as placeholder, not pre-filled

		m.inputs[i] = input
		m.inputLabels[i] = v.Label
	}

	// Skills Folder input (second to last)
	skillsFolderInput := textinput.New()
	skillsFolderInput.Placeholder = "/path/to/.claude/skills/"
	skillsFolderInput.CharLimit = 200
	skillsFolderInput.Width = 50
	if m.skillsFolder != "" {
		skillsFolderInput.SetValue(m.skillsFolder)
	}
	m.inputs[numInputs-2] = skillsFolderInput
	m.inputLabels[numInputs-2] = "Skills Folder"

	// Skill Name input (last) - pre-filled with skill name
	skillNameInput := textinput.New()
	skillNameInput.Placeholder = m.selectedSkill.Name
	skillNameInput.CharLimit = 100
	skillNameInput.Width = 50
	if m.skillFolderName != "" {
		skillNameInput.SetValue(m.skillFolderName)
	} else {
		skillNameInput.SetValue(m.selectedSkill.Name)
	}
	m.inputs[numInputs-1] = skillNameInput
	m.inputLabels[numInputs-1] = "Skill Name"

	// Focus first input
	m.inputFocus = 0
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
	m.inputs[0].Focus()
}

func (m *Model) validateInputs() bool {
	if m.selectedSkill == nil {
		m.errorMsg = "No skill selected"
		return false
	}

	// Validate required variables
	for i, v := range m.selectedSkill.Variables {
		if v.Required && m.inputs[i].Value() == "" {
			m.errorMsg = v.Label + " is required"
			return false
		}
	}

	// Validate skills folder (second to last input)
	skillsFolder := m.inputs[len(m.inputs)-2].Value()
	if skillsFolder == "" {
		m.errorMsg = "Skills Folder is required"
		return false
	}

	// Validate skill name (last input)
	skillName := m.inputs[len(m.inputs)-1].Value()
	if skillName == "" {
		m.errorMsg = "Skill Name is required"
		return false
	}

	m.errorMsg = ""
	return true
}

func (m *Model) saveConfigValues() {
	if m.selectedSkill == nil {
		return
	}

	// Save variable values
	for i, v := range m.selectedSkill.Variables {
		m.configValues[v.Name] = m.inputs[i].Value()
	}

	// Save deploy configuration
	m.skillsFolder = m.inputs[len(m.inputs)-2].Value()
	m.skillFolderName = m.inputs[len(m.inputs)-1].Value()
}

// getDeployPath returns the full deploy path (skillsFolder + skillFolderName)
func (m Model) getDeployPath() string {
	return filepath.Join(m.skillsFolder, m.skillFolderName)
}

func (m Model) updateInputs(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// GetProjectRoot returns the project root, finding it if needed
func GetProjectRoot() string {
	// Try to find project root by looking for go.mod
	dir, _ := os.Getwd()
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// Check if this is SkillFactory
			if _, err := os.Stat(filepath.Join(dir, "skills")); err == nil {
				return dir
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback: assume we're in the right place
	dir, _ = os.Getwd()
	return dir
}
