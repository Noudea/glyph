# glyph

`glyph` is a terminal workspace app for quick tasks and markdown notes.

## Install

### Recommended (binary installer, no clone)

macOS / Linux:

```bash
curl -fsSL https://raw.githubusercontent.com/Noudea/glyph/main/scripts/install.sh | sh
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/Noudea/glyph/main/scripts/install.ps1 | iex
```

Install a specific version (all platforms):

```bash
curl -fsSL https://raw.githubusercontent.com/Noudea/glyph/main/scripts/install.sh | sh -s -- v0.1.0
```

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/Noudea/glyph/main/scripts/install.ps1))) -Version v0.1.0
```

### With Go

```bash
go install github.com/Noudea/glyph/cmd/glyph@latest
glyph
```

### Run from source

```bash
go run ./cmd/glyph
```

## Usage

- Open command palette: `ctrl+p`, `ctrl+k`, `alt+p`
- Next tab: `tab`
- Toggle workspace: `ctrl+w`
- Quit: `ctrl+c`

## Workspaces

- `project`: used when a `.glyph` folder exists in the current directory or an ancestor.
- `global`: fallback to `~/.glyph`.

Data is stored under the active workspace root:

- Tasks: `<workspaceRoot>/tasks/tasks.json`
- Scratchpad: `<workspaceRoot>/scratchpad/scratchpad.md`

Global shortcuts file:

- `~/.glyph/settings/shortcuts.json` (created automatically if missing)

## Requirements

- Go `1.25+`

## Release (maintainers)

Push a version tag:

```bash
git tag v0.1.0
git push origin v0.1.0
```

GitHub Actions publishes release binaries from `.github/workflows/release.yml`.

Local asset build:

```bash
./scripts/build-release-assets.sh v0.1.0 dist
```
