# 🧿 Glyph v0.1 — Product Requirements Document (PRD)

## 1. Overview

**Product Name:** Glyph  
**Version:** v0.1 (MVP)  
**Platform:** Terminal (TUI)  
**Tech Stack:** Go + Charm (Bubble Tea, Lip Gloss, Bubbles, Glamour)

### Vision

Glyph is a playful yet powerful developer workspace inside the terminal.

It helps developers:

- Track active tasks
- Write and preview notes
- Maintain focus within project contexts

Glyph is not a task manager.  
Glyph is not a note-taking app.  
Glyph is a developer focus environment.

---

## 2. Problem Statement

Developers often:

- Switch between multiple projects
- Lose track of active tasks
- Scatter notes across files and apps
- Break flow when switching tools

Existing tools are either:

- Too heavy (Notion, Jira)
- Too simple (todo.txt)
- Not terminal-native
- Not context-aware

Glyph solves this by creating a structured, lightweight, terminal-native workspace for focused building.

---

## 3. Goals (MVP)

1. Provide structured workspaces for context separation
2. Allow full in-TUI editing of notes and tasks
3. Maintain a clean Charm-inspired aesthetic
4. Be fast and keyboard-first
5. Be usable daily by the creator

---

## 4. Non-Goals (v0.1)

- Plugin system
- Command launcher
- Git integration
- Sync/cloud features
- AI features
- Multi-user support
- Due dates, priorities, tagging
- Mobile/GUI version

---

## 5. Core Concepts

### 5.1 Workspace

A workspace is a container for:

- Tasks
- Notes
- Scratchpad

Each workspace is isolated.

Storage structure:

~/.glyph/workspaces/<workspace_name>/
tasks.json
notes.md
scratch.md

---

### 5.2 Apps Inside Glyph

Inside a workspace, Glyph behaves like a small OS with internal “apps”:

- Tasks
- Notes
- Scratch

Each app has:

- Browse mode
- Edit mode
- Preview mode (for markdown content)

---

## 6. Functional Requirements

---

## 6.1 Workspace Management

### Required

- Create workspace
- List workspaces
- Switch workspace
- Persist last opened workspace
- Delete workspace (optional for MVP)

### Behavior

- On first launch → prompt to create workspace
- On subsequent launches → open last used workspace

---

## 6.2 Tasks App

### Task Data Model

Each task contains:

- id
- title
- done (boolean)
- description (markdown, optional)
- created_at

Stored in tasks.json.

---

### Task Features (Must Have)

- Add task
- Edit task title (inline)
- Toggle done
- Delete task
- Persist to disk

---

### Task Detail View

- Selecting a task opens its detail view
- Right pane shows markdown-rendered description
- Edit mode allows editing description via textarea
- Toggle preview/edit mode

---

## 6.3 Notes App

Each workspace has one primary file:

notes.md

### Notes Features

- Preview mode (rendered markdown using Glamour)
- Edit mode (multiline textarea)
- Scrollable viewport
- Save within TUI
- Cancel edit without saving

---

## 6.4 Scratch App

Each workspace has:

scratch.md

### Scratch Features

- Fast open
- Edit inside TUI
- Optional preview toggle
- Used for quick dumps/logs
- Persist to disk

---

## 7. Interaction Model

Glyph operates with clear modes:

- Browse Mode
- Edit Mode
- Preview Mode

The current mode must be clearly visible in the status bar.

---

### Global Keybindings (MVP)

Global:

- q → Quit
- ? → Help
- tab → Switch focus (if split view)
- 1 → Tasks
- 2 → Notes
- 3 → Scratch
- w → Switch workspace

Tasks:

- a → Add task
- d → Toggle done
- x → Delete task
- e → Edit title
- enter → Open detail view

Notes / Scratch:

- e → Edit
- p → Toggle preview
- ctrl+s → Save
- esc → Cancel edit

---

## 8. UI / UX Requirements

- Dark theme default
- Single accent color
- Soft borders (Lip Gloss)
- Clean spacing
- No heavy ASCII art
- Minimal but expressive help bar
- Smooth transitions between modes

The UI must feel:

- Playful
- Intentional
- Fast
- Modern

---

## 9. Performance Requirements

- Startup time < 150ms
- No noticeable UI lag
- No blocking operations on main UI loop
- File operations async-safe where needed

---

## 10. Architecture (High-Level)

### Modules

- workspace/
- tasks/
- notes/
- scratch/
- ui/
- storage/

State managed via Bubble Tea Model.

Each app:

- Owns internal state
- Renders via View()
- Handles Update() messages

---

## 11. Success Criteria

Glyph v0.1 is successful if:

- It is used daily for at least 2 weeks
- Workspaces prevent task/note mixing
- In-TUI editing feels natural
- No desire to immediately replace with external editor
- No performance frustrations

---

## 12. Future Expansion (Post-MVP)

- Command launcher overlay
- Plugin system
- Git awareness
- Multiple notes per workspace
- Session logs
- Task filtering/search
- Command invocation engine

---

# Definition of Done

Glyph v0.1 is complete when:

- Workspaces exist
- Tasks work
- Notes edit + preview works
- Scratch works
- Navigation feels clean
- It is actually usable for real projects
