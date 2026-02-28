package shell

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/marketplace"
)

type configFile struct {
	Version   int                        `json:"version,omitempty"`
	Commands  []commandConfig            `json:"commands"`
	Shortcuts map[string]json.RawMessage `json:"shortcuts"`
}

type commandConfig struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Run     string `json:"run,omitempty"`
	Script  string `json:"script,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type configWriteFile struct {
	Version   int                 `json:"version,omitempty"`
	Commands  []commandConfig     `json:"commands"`
	Shortcuts map[string][]string `json:"shortcuts"`
}

func configPath(globalRoot string) string {
	return filepath.Join(globalRoot, "settings", "config.json")
}

func defaultConfigTemplate() configWriteFile {
	return configWriteFile{
		Version:  1,
		Commands: []commandConfig{},
		Shortcuts: map[string][]string{
			commandLauncherOpen: {"ctrl+p", "ctrl+k", "alt+p"},
		},
	}
}

func ensureGlobalConfig(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return writeDefaultConfig(path)
}

func writeDefaultConfig(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	template := defaultConfigTemplate()
	template.Shortcuts = normalizeShortcutMap(template.Shortcuts)
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func loadConfig(path string) (configFile, error) {
	var out configFile
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, err
	}
	if out.Shortcuts == nil {
		out.Shortcuts = map[string]json.RawMessage{}
	}
	return out, nil
}

func findNearestProjectConfig(startDir string) (string, bool, error) {
	if strings.TrimSpace(startDir) == "" {
		return "", false, nil
	}
	current := filepath.Clean(startDir)
	for {
		path := filepath.Join(current, ".glyph", "config.json")
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			return path, true, nil
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return "", false, err
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return "", false, nil
}

func decodeShortcutMap(input map[string]json.RawMessage) (map[string][]string, []error) {
	out := make(map[string][]string, len(input))
	var errs []error
	for commandID, raw := range input {
		var list []string
		if err := json.Unmarshal(raw, &list); err == nil {
			out[commandID] = list
			continue
		}
		var single string
		if err := json.Unmarshal(raw, &single); err == nil {
			out[commandID] = []string{single}
			continue
		}
		errs = append(errs, errors.New("invalid shortcut format for "+commandID))
	}
	return out, errs
}

func mergeCommands(global []commandConfig, globalRoot string, project []commandConfig, projectRoot string) ([]core.Command, []error) {
	commandsByID := make(map[string]core.Command)
	order := make([]string, 0, len(global)+len(project))
	orderSet := make(map[string]struct{})
	var errs []error

	apply := func(entries []commandConfig, source string, configRoot string) {
		for _, item := range entries {
			command, ok, err := parseCommandConfig(item, source, configRoot)
			if err != nil {
				errs = append(errs, err)
			}
			if !ok {
				continue
			}
			if _, exists := orderSet[command.ID]; !exists {
				orderSet[command.ID] = struct{}{}
				order = append(order, command.ID)
			}
			commandsByID[command.ID] = command
		}
	}

	apply(global, commandSourceGlobal, globalRoot)
	apply(project, commandSourceProject, projectRoot)

	out := make([]core.Command, 0, len(order))
	for _, id := range order {
		if command, ok := commandsByID[id]; ok {
			out = append(out, command)
		}
	}

	sort.SliceStable(out, func(i, j int) bool {
		left := strings.ToLower(strings.TrimSpace(out[i].Label))
		right := strings.ToLower(strings.TrimSpace(out[j].Label))
		if left == right {
			return out[i].ID < out[j].ID
		}
		return left < right
	})

	return out, errs
}

func parseCommandConfig(item commandConfig, source string, configRoot string) (core.Command, bool, error) {
	id := strings.TrimSpace(item.ID)
	if id == "" {
		return core.Command{}, false, errors.New("command id is required")
	}
	if !commandEnabled(item.Enabled) {
		return core.Command{}, false, nil
	}

	run := strings.TrimSpace(item.Run)
	script := strings.TrimSpace(item.Script)

	if run == "" && script == "" {
		return core.Command{}, false, errors.New("command run or script is required for " + id)
	}

	// script is resolved to an absolute path relative to configRoot.
	if script != "" {
		run = filepath.Join(configRoot, script)
	}

	label := strings.TrimSpace(item.Label)
	if label == "" {
		label = id
	}

	return core.Command{
		ID:      id,
		Label:   label,
		Kind:    core.CommandExec,
		Group:   "commands",
		Run:     run,
		Source:  source,
		Managed: source == commandSourceManaged,
	}, true, nil
}

func commandEnabled(flag *bool) bool {
	if flag == nil {
		return true
	}
	return *flag
}

func normalizeShortcutMap(input map[string][]string) map[string][]string {
	out := make(map[string][]string, len(input))
	for commandID, keys := range input {
		out[commandID] = normalizeShortcutKeys(keys)
	}
	return out
}

func loadSpellbookCommands(root string, source string, absolute bool) ([]core.Command, error) {
	installed := marketplace.ListInstalled(root)
	var commands []core.Command
	for id, sb := range installed {
		for _, cmd := range sb.Commands {
			if cmd.Enabled != nil && !*cmd.Enabled {
				continue
			}

			// Resolve the script/run to an executable path.
			// Spellbooks should use "script" (file relative to spellbook dir),
			// but we also support "run" for backwards compatibility.
			scriptFile := cmd.Script
			if scriptFile == "" {
				scriptFile = cmd.Run
			}

			var run string
			if scriptFile != "" {
				if absolute {
					run = filepath.Join(root, "spellbooks", id, scriptFile)
				} else {
					run = filepath.Join(".glyph", "spellbooks", id, scriptFile)
				}
			} else {
				run = cmd.Run
			}

			label := strings.TrimSpace(cmd.Label)
			if label == "" {
				label = cmd.ID
			}
			commands = append(commands, core.Command{
				ID:     cmd.ID,
				Label:  label,
				Kind:   core.CommandExec,
				Group:  "spellbook",
				Run:    run,
				Source: source,
			})
		}
	}
	return commands, nil
}

func (m *Model) reloadConfig() error {
	var problems []error

	globalRoot, err := m.resolver.ResolveGlobal()
	if err != nil {
		return err
	}

	m.globalConfigPath = configPath(globalRoot.RootPath)
	if err := ensureGlobalConfig(m.globalConfigPath); err != nil {
		return err
	}

	globalConfig, err := loadConfig(m.globalConfigPath)
	if err != nil {
		problems = append(problems, err)
		globalConfig = configFile{
			Version:   1,
			Commands:  []commandConfig{},
			Shortcuts: map[string]json.RawMessage{},
		}
	}

	projectPath, found, err := findNearestProjectConfig(m.startDir)
	if err != nil {
		problems = append(problems, err)
	}

	projectConfig := configFile{
		Version:   1,
		Commands:  []commandConfig{},
		Shortcuts: map[string]json.RawMessage{},
	}
	if found {
		m.projectConfigPath = projectPath
		projectLoaded, loadErr := loadConfig(projectPath)
		if loadErr != nil {
			problems = append(problems, loadErr)
		} else {
			projectConfig = projectLoaded
		}
	} else {
		m.projectConfigPath = ""
	}

	projectRoot := ""
	if m.projectConfigPath != "" {
		projectRoot = filepath.Dir(m.projectConfigPath) // .glyph/ directory
	}
	commands, commandProblems := mergeCommands(globalConfig.Commands, globalRoot.RootPath, projectConfig.Commands, projectRoot)
	problems = append(problems, commandProblems...)

	// Load spellbook commands from installed spellbooks.
	globalSpellbookCmds, err := loadSpellbookCommands(globalRoot.RootPath, commandSourceSpellbook, true)
	if err != nil {
		problems = append(problems, err)
	}
	commands = append(commands, globalSpellbookCmds...)

	if projectRoot != "" {
		projectSpellbookCmds, err := loadSpellbookCommands(projectRoot, commandSourceSpellbook, false)
		if err != nil {
			problems = append(problems, err)
		}
		commands = append(commands, projectSpellbookCmds...)
	}

	if m.state != nil {
		m.state.Commands = commands
	}

	shortcutOverrides, shortcutProblems := decodeShortcutMap(globalConfig.Shortcuts)
	problems = append(problems, shortcutProblems...)
	if err := m.applyShortcuts(shortcutOverrides); err != nil {
		problems = append(problems, err)
	}

	return errors.Join(problems...)
}
