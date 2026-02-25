package shell

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type shortcutsConfig struct {
	Version   int                        `json:"version,omitempty"`
	Shortcuts map[string]json.RawMessage `json:"shortcuts"`
}

type shortcutsWriteConfig struct {
	Version   int                 `json:"version,omitempty"`
	Shortcuts map[string][]string `json:"shortcuts"`
}

func shortcutsConfigPath(globalRoot string) string {
	return filepath.Join(globalRoot, "settings", "shortcuts.json")
}

func loadShortcutsConfig(path string) (map[string][]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var config shortcutsConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	out := make(map[string][]string, len(config.Shortcuts))
	for commandID, raw := range config.Shortcuts {
		var list []string
		if err := json.Unmarshal(raw, &list); err == nil {
			out[commandID] = list
			continue
		}
		var single string
		if err := json.Unmarshal(raw, &single); err == nil {
			out[commandID] = []string{single}
		}
	}
	return out, nil
}

func ensureShortcutsConfig(path string, defaults map[string][]string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeShortcutsConfig(path, defaults)
}

func writeShortcutsConfig(path string, shortcuts map[string][]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	config := shortcutsWriteConfig{
		Version:   1,
		Shortcuts: make(map[string][]string, len(shortcuts)),
	}
	for commandID, keys := range shortcuts {
		config.Shortcuts[commandID] = normalizeShortcutKeys(keys)
	}
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
