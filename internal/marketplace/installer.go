package marketplace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Install downloads a community spellbook and installs it under root/spellbooks/<id>/.
func Install(root, id string) error {
	registry, err := FetchRegistry()
	if err != nil {
		return err
	}
	sb, ok := registry[id]
	if !ok {
		return fmt.Errorf("marketplace: spellbook %q not found in registry", id)
	}

	dir := filepath.Join(root, "spellbooks", id)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("marketplace: create dir: %w", err)
	}

	for _, cmd := range sb.Commands {
		filename := cmd.Script
		if filename == "" {
			filename = cmd.Run
		}
		if filename == "" {
			continue
		}
		data, err := FetchScript(id, filename)
		if err != nil {
			return err
		}
		dest := filepath.Join(dir, filename)
		if err := os.WriteFile(dest, data, 0o755); err != nil {
			return fmt.Errorf("marketplace: write script %s: %w", filename, err)
		}
	}

	manifest, err := json.MarshalIndent(sb, "", "  ")
	if err != nil {
		return fmt.Errorf("marketplace: marshal manifest: %w", err)
	}
	manifest = append(manifest, '\n')
	if err := os.WriteFile(filepath.Join(dir, "spellbook.json"), manifest, 0o644); err != nil {
		return fmt.Errorf("marketplace: write manifest: %w", err)
	}

	return nil
}

// Uninstall removes an installed spellbook.
func Uninstall(root, id string) error {
	dir := filepath.Join(root, "spellbooks", id)
	return os.RemoveAll(dir)
}

// ListInstalled scans root/spellbooks/ and returns all installed spellbooks keyed by ID.
func ListInstalled(root string) map[string]Spellbook {
	out := make(map[string]Spellbook)

	pattern := filepath.Join(root, "spellbooks", "*", "spellbook.json")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return out
	}

	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var sb Spellbook
		if err := json.Unmarshal(data, &sb); err != nil {
			continue
		}
		id := filepath.Base(filepath.Dir(path))
		out[id] = sb
	}
	return out
}

// NeedsUpdate returns true when the remote spellbook has a newer version.
func NeedsUpdate(local, remote Spellbook) bool {
	return local.Version != remote.Version
}

// Update re-installs a spellbook by removing the old version and installing the new one.
func Update(root, id string) error {
	if err := Uninstall(root, id); err != nil {
		return err
	}
	return Install(root, id)
}
