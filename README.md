# glyph

`glyph` is a terminal command palette launcher for CLI/TUI workflows.

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
- Quit: `ctrl+c`

Commands run in:

- Current terminal session
- Current folder where `glyph` started

## Configuration

Global config file (auto-created if missing):

- `~/.glyph/settings/config.json`

Optional project config (nearest ancestor):

- `<repo>/.glyph/config.json`

Project config can add/override `commands` by `id`.  
Project `shortcuts` are ignored in v1.

### Config schema

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
      "id": "user.ls",
      "label": "List Files",
      "run": "ls -la",
      "enabled": true
    }
  ],
  "shortcuts": {
    "launcher.open": ["ctrl+p", "ctrl+k", "alt+p"],
    "user.lazygit": ["ctrl+g"]
  }
}
```

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
