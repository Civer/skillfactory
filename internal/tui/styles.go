// Package tui provides the terminal user interface for SkillFactory
package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7C3AED") // Purple
	accentColor    = lipgloss.Color("#A78BFA") // Light Purple
	highlightColor = lipgloss.Color("#C4B5FD") // Even lighter purple
	secondaryColor = lipgloss.Color("#10B981") // Green
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	errorColor     = lipgloss.Color("#EF4444") // Red
	successColor   = lipgloss.Color("#10B981") // Green
	darkBgColor    = lipgloss.Color("#1E1B4B") // Dark purple background

	// Header styles (like Claude Code)
	logoStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	// Content styles
	selectedStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	mutedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	successStyle = lipgloss.NewStyle().
			Foreground(successColor)

	errorStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	inputLabelStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	versionStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Faint(true)
)
