package shell

import (
	"sort"
	"strings"
)

const commandLauncherOpen = "launcher.open"

var defaultMainCommandShortcuts = map[string][]string{
	commandLauncherOpen:   {"ctrl+p", "ctrl+k", "alt+p"},
	actionTabsNext:        {"tab"},
	actionWorkspaceToggle: {"ctrl+w"},
}

func (m *Model) loadGlobalShortcuts() error {
	global, err := m.resolver.ResolveGlobal()
	if err != nil {
		return err
	}
	path := shortcutsConfigPath(global.RootPath)
	if err := ensureShortcutsConfig(path, defaultMainCommandShortcuts); err != nil {
		return err
	}
	overrides, err := loadShortcutsConfig(path)
	if err != nil {
		return err
	}
	m.applyShortcuts(overrides)
	return nil
}

func (m *Model) resolveMainCommandIDForKey(key string) (string, bool) {
	if m.shortcutCommands == nil {
		m.applyShortcuts(nil)
	}
	commandID, ok := m.shortcutCommands[normalizeShortcutKey(key)]
	return commandID, ok
}

func (m Model) primaryShortcut(commandID string, fallback string) string {
	shortcuts := m.commandShortcuts[commandID]
	if len(shortcuts) == 0 {
		return fallback
	}
	return shortcuts[0]
}

func (m Model) shortcutsHint(commandID string, fallback string) string {
	shortcuts := m.commandShortcuts[commandID]
	if len(shortcuts) == 0 {
		return fallback
	}
	return strings.Join(shortcuts, "/")
}

func (m *Model) applyShortcuts(overrides map[string][]string) {
	bindings := copyShortcutBindings(defaultMainCommandShortcuts)
	known := m.knownCommandIDs()
	commandIDs := make([]string, 0, len(overrides))
	for commandID := range overrides {
		commandIDs = append(commandIDs, commandID)
	}
	sort.Strings(commandIDs)
	for _, commandID := range commandIDs {
		if _, ok := known[commandID]; !ok {
			continue
		}
		keys := normalizeShortcutKeys(overrides[commandID])
		bindings[commandID] = keys
	}

	reverse := make(map[string]string)
	allIDs := make([]string, 0, len(bindings))
	for commandID := range bindings {
		allIDs = append(allIDs, commandID)
	}
	sort.Strings(allIDs)
	for _, commandID := range allIDs {
		for _, key := range bindings[commandID] {
			reverse[key] = commandID
		}
	}
	m.commandShortcuts = bindings
	m.shortcutCommands = reverse
}

func copyShortcutBindings(input map[string][]string) map[string][]string {
	out := make(map[string][]string, len(input))
	for commandID, keys := range input {
		copied := make([]string, 0, len(keys))
		for _, key := range keys {
			key = normalizeShortcutKey(key)
			if key != "" {
				copied = append(copied, key)
			}
		}
		out[commandID] = copied
	}
	return out
}

func normalizeShortcutKeys(keys []string) []string {
	seen := make(map[string]struct{}, len(keys))
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		key = normalizeShortcutKey(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, key)
	}
	return out
}

func normalizeShortcutKey(key string) string {
	return strings.ToLower(strings.TrimSpace(key))
}

func (m Model) knownCommandIDs() map[string]struct{} {
	ids := map[string]struct{}{
		commandLauncherOpen:    {},
		actionTabsNext:         {},
		actionWorkspaceToggle:  {},
		actionWorkspaceProject: {},
		actionWorkspaceCreate:  {},
		actionWorkspaceGlobal:  {},
	}
	if m.state != nil {
		for _, cmd := range m.state.Commands {
			if strings.TrimSpace(cmd.ID) == "" {
				continue
			}
			ids[cmd.ID] = struct{}{}
		}
	}
	return ids
}
