# Vikunja Skill

**USE THIS SKILL** when the user asks to:
- Create, update, delete, or list tasks/todos/Aufgaben
- Manage projects or labels in Vikunja
- Mark tasks as done
- Any task management operations

This skill provides direct CLI access to Vikunja - prefer this over MCP for lower token usage.

Base directory: {{SKILL_PATH}}

## Overview

Task management via Vikunja API with lean JSON output (~95% fewer tokens than MCP integration).

## Project IDs

Reference table for commonly used projects:

| ID | Name | Description |
|----|------|-------------|
| 1 | Inbox | Default project for quick captures |
| 2 | Work | Work-related tasks |

## Label IDs

Reference table for commonly used labels:

| ID | Name |
|----|------|
| 1 | urgent |
| 2 | waiting |

## Notes

- All responses are lean JSON with minimal overhead
- Date format for `--due`, `--start`, `--end`: `YYYY-MM-DD`
- Priority: 0 (none) to 5 (highest)
- Labels: Use `--labels 1,2,3` on create, or `add-label`/`remove-label` commands

## Commands

{{COMMANDS}}
