# Glyph v1 - Product Requirements Document (PRD)

## 1. Overview

**Product Name:** Glyph  
**Version:** v1 (MVP pivot)  
**Platform:** Terminal (TUI)  
**Tech Stack:** Go + Charm (Bubble Tea, Lip Gloss, Bubbles)

### Product Direction

Glyph is a super command palette launcher for terminal workflows.

It helps developers:

- Launch frequently used CLI/TUI tools fast
- Trigger repeatable shell commands from keyboard shortcuts
- Stay in flow without leaving the terminal

Glyph v1 is not a task manager and not a notes app.

## 2. Problem Statement

Terminal users repeatedly type the same commands and switch between tools manually.
This creates friction, context switching, and inconsistent workflows.

Glyph solves this with a keyboard-first command palette plus user-defined shortcuts.

## 3. Goals (v1)

1. Launch commands quickly from a searchable palette.
2. Let users and teams define custom commands, and let users define personal shortcuts.
3. Execute commands in the current terminal and current folder.
4. Keep startup and interaction fast enough for daily use.
5. Provide predictable behavior and clear failure feedback.

## 4. Non-Goals (v1)

- Marketplace / package installation
- Launching in a new terminal window
- Background jobs / daemonized execution
- Project-level shortcuts (team keymaps) in v1
- Plugin system, sync/cloud, AI features

## 5. Core Product Principles

1. Keyboard-first: open, search, run with minimal keystrokes.
2. Predictable execution: always run where Glyph started.
3. Minimal config surface: simple JSON, low ceremony.
4. Fail safely: invalid config or command failures must not crash the app.

## 6. Functional Requirements

### 6.1 Command Palette

Required behavior:

- Open palette via global shortcut(s).
- Filter commands by label, ID, and shortcut text.
- Navigate with keyboard and run selected command with `enter`.
- Close palette with `esc`.

### 6.2 Custom Commands

Users can define commands in a config file.

Each command must support:

- `id` (string, unique)
- `label` (string, human-readable)
- `run` (string, shell command)
- `enabled` (boolean, optional; default `true`)

Rules:

- Disabled commands are hidden and cannot be executed.
- Invalid entries are skipped with a visible error hint.
- Duplicate IDs are rejected deterministically (first valid wins).

### 6.3 Custom Shortcuts

Users can map shortcuts to command IDs.

Rules:

- Shortcuts are normalized and case-insensitive.
- Shortcut collisions are rejected with clear error messaging.
- Reserved keys cannot be overridden (at minimum app quit and launcher toggle keys).
- A command can have multiple shortcuts; first is primary display shortcut.

### 6.4 Execution Model

When a command is launched (palette or shortcut):

- Execute in the current terminal session (no new window).
- Execute in the current working directory from which Glyph was started.
- Hand terminal control to the launched process.
- Return cleanly to Glyph when process exits.

### 6.5 Config Resolution and Merge

On startup, Glyph resolves configuration in this order:

- Load global config from `~/.glyph/settings/config.json`.
- Detect nearest ancestor project config at `.glyph/config.json` from the startup folder.
- If found, auto-load project config with no confirmation prompt.

Merge rules:

- `commands`: merge by `id`, where project command definitions override global definitions.
- `shortcuts`: load from global config only; project shortcuts are ignored in v1.

### 6.6 Error Handling

Must handle gracefully:

- Invalid/missing global config file
- Invalid project config file
- Unknown shortcut target command ID
- Command executable not found
- Non-zero command exit codes

Failures should surface a concise error in the hint/status area.

## 7. Configuration and Storage

### 7.1 Files

- Global configuration: `~/.glyph/settings/config.json`
- Optional project configuration: `<repo>/.glyph/config.json` (nearest ancestor from startup folder)

### 7.2 config.json (draft)

```json
{
  "version": 1,
  "commands": [
    {
      "id": "user.lazygit",
      "label": "LazyGit",
      "run": "lazygit",
      "enabled": true
    },
    {
      "id": "user.test",
      "label": "Run Tests",
      "run": "npm test",
      "enabled": true
    }
  ],
  "shortcuts": {
    "launcher.open": [
      "ctrl+p",
      "ctrl+k",
      "alt+p"
    ],
    "user.lazygit": [
      "ctrl+g"
    ],
    "user.test": [
      "ctrl+t"
    ]
  }
}
```

### 7.3 Merge Behavior

- Both files use the same schema.
- Global config is the base layer.
- Project config can add/override `commands`.
- Project `shortcuts` are ignored in v1.

## 8. User Flows

1. Launch from palette:
   Open palette, search command, press `enter`, command runs, Glyph resumes on exit.
2. Launch from shortcut:
   Press mapped shortcut in main mode, command runs immediately, Glyph resumes on exit.
3. Configure personal commands:
   Edit `config.json` (`commands` section), restart/reload Glyph, new commands appear.
4. Configure shared project commands:
   Commit `.glyph/config.json` in a repo; when Glyph starts in that repo, project commands are auto-loaded and merged.
5. Configure personal shortcuts:
   Edit global `config.json` (`shortcuts` section), restart/reload Glyph, shortcuts apply.
6. Failure case:
   Command fails or cannot start, Glyph returns and shows a concise error.

## 9. Performance and Reliability Requirements

- Startup target: < 150ms on typical developer machines
- No noticeable lag while typing/filtering in palette
- No blocking UI behavior beyond intended process handoff during command execution
- No crashes from malformed config files

## 10. Acceptance Criteria

1. Missing global `config.json` is created with a valid template including default launcher shortcuts.
2. Invalid config never crashes Glyph.
3. Mixed valid/invalid command entries load valid entries only.
4. Disabled commands do not appear in palette and cannot run by shortcut.
5. Palette filtering matches `label`, `id`, and shortcut.
6. `enter` on a selected command executes its `run` string.
7. Shortcut-triggered execution matches palette execution behavior.
8. Commands always run in the current terminal session.
9. Commands always run in Glyph's startup working directory.
10. After command exit, Glyph resumes interactive TUI normally.
11. Shortcut collisions are rejected with clear feedback.
12. Reserved shortcuts cannot be overridden.
13. If `.glyph/config.json` exists in the nearest ancestor, project commands are auto-loaded.
14. Project command `id` collisions override global command definitions deterministically.
15. Project shortcuts do not override global shortcuts in v1.

## 11. Future Expansion (Post-v1)

- Marketplace for installing CLI/TUI packages
- Install/update/uninstall command packs from within palette
- Verified publishers and package trust/signature model
- Optional run targets (current terminal vs new terminal window)
- Optional per-command working directory override

## Definition of Done

Glyph v1 is complete when:

1. Users can define and launch custom commands from palette and shortcuts.
2. Execution is reliable in current terminal + current folder.
3. Config and execution failures are handled without crashes.
4. The app is stable enough for daily command-launcher usage.
