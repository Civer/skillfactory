// SkillFactory - Build & Deploy Claude Code Skills
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/petervogelmann/skillfactory/internal/tui"
)

// version is set via ldflags at build time
// e.g.: go build -ldflags "-X main.version=$(git describe --tags --always)"
var version = "dev"

func main() {
	// Handle --version flag
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("SkillFactory", version)
		os.Exit(0)
	}

	// Find project root
	projectRoot := tui.GetProjectRoot()

	// Create and run TUI
	model := tui.NewModel(projectRoot, version)
	p := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
