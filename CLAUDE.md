# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Git Commits

When creating commits, do NOT include:
- "🤖 Generated with Claude Code" footer
- "Co-Authored-By: Claude" lines

Keep commit messages clean and focused on the changes.

## Project Overview

SkillFactory is a Go-based tool for building and deploying "Skills" for Claude Code. Skills are CLI tools that output lean JSON, optimized for minimal token usage compared to MCP servers.

## Build & Run Commands

```bash
# Build the main TUI application (with version from git tag)
go build -ldflags "-X main.version=$(git describe --tags --always)" -o skillfactory ./cmd/skillfactory

# Build without version (shows "dev")
go build -o skillfactory ./cmd/skillfactory

# Build a skill directly (e.g., vikunja)
cd skills/vikunja && go build -o vikunja .

# Run the TUI to configure and deploy skills
./skillfactory

# Check version
./skillfactory --version
```

## Architecture

### Core Components

- **cmd/skillfactory/main.go** - Entry point for the TUI application
- **internal/tui/** - Bubbletea-based TUI with multi-view navigation
  - `model.go` - State management, views (SkillList → Config → Confirm → Overwrite → Building → Done)
  - `view.go` - Rendering functions for each view
  - `commands.go` - Build, deploy, and documentation generation logic
  - `styles.go` - Lipgloss styling
- **internal/skill/manifest.go** - Parses `skill.yaml` manifests, discovers skills, handles SkillErrors

### Skill Structure

Each skill lives in `skills/<name>/` and requires:
- `skill.yaml` - Manifest defining name, variables, build config, and deploy settings
- `main.go` - Cobra-based CLI entry point
- Domain packages organized by entity (e.g., `client/`, `tasks/`, `labels/`, `projects/`)

### Vikunja Skill Structure (Example)

```
skills/vikunja/
├── skill.yaml           # Manifest with variables (VIKUNJA_URL, VIKUNJA_TOKEN, PROJECT_IDS)
├── main.go              # Entry point, loads .env, registers commands
├── client/
│   └── client.go        # HTTP client with auth, Get/Post/Put/Delete methods
├── tasks/
│   ├── types.go         # Task, TaskLean, CreateTaskRequest, UpdateTaskRequest, Label
│   ├── service.go       # CRUD + Label operations (GetLabels, AddLabel, RemoveLabel)
│   └── commands.go      # Cobra commands: list, get, create, update, delete, done, labels, add-label, remove-label
├── labels/
│   ├── types.go         # Label, LabelLean
│   ├── service.go       # CRUD operations
│   └── commands.go      # Cobra commands
├── projects/
│   ├── types.go         # Project, ProjectLean
│   ├── service.go       # CRUD operations
│   └── commands.go      # Cobra commands
└── SKILL.template.md    # Template for SKILL.md generation
```

### Data Flow

1. TUI discovers skills by scanning `skills/*/skill.yaml`
2. User configures environment variables defined in the manifest
3. Deploys to: `<skills-folder>/<skill-name>/` (two separate inputs)
4. Build: `go build` compiles skill to `dist/<binary>`
5. Deploy:
   - Copies binary to `bin/<binary>`
   - Generates `.env` file with configured variables (loaded via godotenv)
   - Generates `SKILL.md` with YAML frontmatter (name, description) for Claude Code discovery

### Key Patterns

- **Lean JSON Output**: Types have `ToLean()` methods returning minimal structs for reduced token usage
- **Environment Variables**: Loaded via godotenv from `.env` file in binary directory
- **Date Formatting**: `formatDate()` in service.go converts `YYYY-MM-DD` to RFC3339 with local timezone
- **TUI Architecture**: Bubbletea's Elm pattern (Model → Update → View)
- **skill.yaml Variables**: Types `string`, `secret` (masked), `json`
- **Labels**: Read-only on Task, managed via separate endpoints (add-label, remove-label)
- **Overwrite Warning**: TUI checks if skill exists before deploying

### API Integration (Vikunja Example)

- Tasks support: title, description, priority, due_date, start_date, end_date, hex_color, is_favorite, percent_done
- Labels are attached via `PUT /tasks/{id}/labels` with `{label_id: N}`
- Labels removed via `DELETE /tasks/{id}/labels/{label_id}`
- Dates: API expects RFC3339 (`2025-12-10T00:00:00+01:00`), CLI accepts `YYYY-MM-DD`

## Common Tasks

### Adding a New Flag to a Command

1. Add field to request struct in `types.go`
2. Add flag variable and `Flags().XxxVar()` call in `commands.go`
3. Set field in request if flag changed
4. If date field: ensure `formatDate()` is called in service

### Adding a New Subcommand

1. Create the `*cobra.Command` in `commands.go`
2. Add service method in `service.go` if needed
3. Add to `cmd.AddCommand(...)` list

### Testing a Skill

```bash
cd skills/vikunja
go build -o vikunja .
./vikunja tasks create --title "Test" --project 2 --due 2025-12-15 --labels 1,3
./vikunja tasks labels 123
./vikunja tasks add-label 123 5
```
