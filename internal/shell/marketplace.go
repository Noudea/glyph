package shell

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/Noudea/glyph/internal/marketplace"
	tea "github.com/charmbracelet/bubbletea"
)

const commandMarketplaceOpen = "marketplace.open"

// Messages for async marketplace operations.
type marketplaceListMsg struct {
	entries []marketplaceEntry
	err     string
}

type marketplaceInstallMsg struct {
	id  string
	err string
}

type marketplaceUninstallMsg struct {
	id  string
	err string
}

type marketplaceUpdateMsg struct {
	id  string
	err string
}

func (m *Model) openMarketplace() tea.Cmd {
	m.mode = ModeMarketplace
	m.marketplace = marketplaceState{
		loading: true,
	}
	return m.fetchMarketplaceList()
}

func (m *Model) fetchMarketplaceList() tea.Cmd {
	return func() tea.Msg {
		registry, err := marketplace.FetchRegistry()
		if err != nil {
			return marketplaceListMsg{err: err.Error()}
		}

		var entries []marketplaceEntry
		for id, sb := range registry {
			entries = append(entries, marketplaceEntry{
				ID:     id,
				Remote: sb,
			})
		}

		// Check global installed.
		globalRoot, _ := m.resolveGlobalRoot()
		globalInstalled := marketplace.ListInstalled(globalRoot)

		// Check project installed (if in a project).
		var projectInstalled map[string]marketplace.Spellbook
		projectRoot := m.resolveProjectRoot()
		if projectRoot != "" {
			projectInstalled = marketplace.ListInstalled(projectRoot)
		}

		for i, e := range entries {
			if local, ok := globalInstalled[e.ID]; ok {
				entries[i].InstalledGlobal = true
				entries[i].HasUpdateGlobal = marketplace.NeedsUpdate(local, e.Remote)
			}
			if local, ok := projectInstalled[e.ID]; ok {
				entries[i].InstalledProject = true
				entries[i].HasUpdateProject = marketplace.NeedsUpdate(local, e.Remote)
			}
		}

		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Remote.Name) < strings.ToLower(entries[j].Remote.Name)
		})

		return marketplaceListMsg{entries: entries}
	}
}

func (m *Model) marketplaceInstallGlobal(id string) tea.Cmd {
	m.marketplace.installing = id
	return func() tea.Msg {
		root, err := m.resolveGlobalRoot()
		if err != nil {
			return marketplaceInstallMsg{id: id, err: err.Error()}
		}
		if err := marketplace.Install(root, id); err != nil {
			return marketplaceInstallMsg{id: id, err: err.Error()}
		}
		return marketplaceInstallMsg{id: id}
	}
}

func (m *Model) marketplaceInstallProject(id string) tea.Cmd {
	m.marketplace.installing = id
	return func() tea.Msg {
		root := m.resolveProjectRoot()
		if root == "" {
			root = filepath.Join(m.startDir, ".glyph")
		}
		if err := marketplace.Install(root, id); err != nil {
			return marketplaceInstallMsg{id: id, err: err.Error()}
		}
		return marketplaceInstallMsg{id: id}
	}
}

func (m *Model) marketplaceUninstall(id string) tea.Cmd {
	return func() tea.Msg {
		// Uninstall from both scopes where installed.
		globalRoot, _ := m.resolveGlobalRoot()
		globalInstalled := marketplace.ListInstalled(globalRoot)
		if _, ok := globalInstalled[id]; ok {
			if err := marketplace.Uninstall(globalRoot, id); err != nil {
				return marketplaceUninstallMsg{id: id, err: err.Error()}
			}
		}

		projectRoot := m.resolveProjectRoot()
		if projectRoot != "" {
			projectInstalled := marketplace.ListInstalled(projectRoot)
			if _, ok := projectInstalled[id]; ok {
				if err := marketplace.Uninstall(projectRoot, id); err != nil {
					return marketplaceUninstallMsg{id: id, err: err.Error()}
				}
			}
		}

		return marketplaceUninstallMsg{id: id}
	}
}

func (m *Model) marketplaceUpdate(id string) tea.Cmd {
	m.marketplace.installing = id
	return func() tea.Msg {
		// Update in all scopes where installed.
		globalRoot, _ := m.resolveGlobalRoot()
		globalInstalled := marketplace.ListInstalled(globalRoot)
		if _, ok := globalInstalled[id]; ok {
			if err := marketplace.Update(globalRoot, id); err != nil {
				return marketplaceUpdateMsg{id: id, err: err.Error()}
			}
		}

		projectRoot := m.resolveProjectRoot()
		if projectRoot != "" {
			projectInstalled := marketplace.ListInstalled(projectRoot)
			if _, ok := projectInstalled[id]; ok {
				if err := marketplace.Update(projectRoot, id); err != nil {
					return marketplaceUpdateMsg{id: id, err: err.Error()}
				}
			}
		}

		return marketplaceUpdateMsg{id: id}
	}
}

func (m *Model) updateMarketplace(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case marketplaceListMsg:
		m.marketplace.loading = false
		if msg.err != "" {
			m.marketplace.err = msg.err
		} else {
			m.marketplace.entries = msg.entries
		}
		return m, nil

	case marketplaceInstallMsg:
		m.marketplace.installing = ""
		if msg.err != "" {
			m.marketplace.err = msg.err
		} else {
			_ = m.reloadConfig()
		}
		return m, m.fetchMarketplaceList()

	case marketplaceUninstallMsg:
		if msg.err != "" {
			m.marketplace.err = msg.err
		} else {
			_ = m.reloadConfig()
		}
		return m, m.fetchMarketplaceList()

	case marketplaceUpdateMsg:
		m.marketplace.installing = ""
		if msg.err != "" {
			m.marketplace.err = msg.err
		} else {
			_ = m.reloadConfig()
		}
		return m, m.fetchMarketplaceList()

	case tea.KeyMsg:
		return m.handleMarketplaceKey(msg)
	}
	return m, nil
}

func (m *Model) handleMarketplaceKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// If we're in the install-scope prompt, handle g/p/esc.
	if m.marketplace.confirmInstall != "" {
		switch key {
		case "g":
			id := m.marketplace.confirmInstall
			m.marketplace.confirmInstall = ""
			return m, m.marketplaceInstallGlobal(id)
		case "p":
			id := m.marketplace.confirmInstall
			m.marketplace.confirmInstall = ""
			return m, m.marketplaceInstallProject(id)
		case "esc":
			m.marketplace.confirmInstall = ""
			return m, nil
		}
		return m, nil
	}

	switch key {
	case "esc":
		m.openLauncher()
		return m, nil

	case "up", "k":
		if m.marketplace.cursor > 0 {
			m.marketplace.cursor--
		}
		return m, nil

	case "down", "j":
		if m.marketplace.cursor < len(m.marketplace.entries)-1 {
			m.marketplace.cursor++
		}
		return m, nil

	case "i":
		entries := m.marketplace.entries
		if len(entries) > 0 && m.marketplace.cursor >= 0 && m.marketplace.cursor < len(entries) {
			e := entries[m.marketplace.cursor]
			if !e.InstalledGlobal || !e.InstalledProject {
				// If in a project, ask scope. Otherwise install globally.
				if m.resolveProjectRoot() != "" {
					m.marketplace.confirmInstall = e.ID
				} else {
					return m, m.marketplaceInstallGlobal(e.ID)
				}
			}
		}
		return m, nil

	case "u":
		entries := m.marketplace.entries
		if len(entries) > 0 && m.marketplace.cursor >= 0 && m.marketplace.cursor < len(entries) {
			e := entries[m.marketplace.cursor]
			if e.InstalledGlobal || e.InstalledProject {
				return m, m.marketplaceUninstall(e.ID)
			}
		}
		return m, nil

	case "U":
		entries := m.marketplace.entries
		if len(entries) > 0 && m.marketplace.cursor >= 0 && m.marketplace.cursor < len(entries) {
			e := entries[m.marketplace.cursor]
			if e.HasUpdateGlobal || e.HasUpdateProject {
				return m, m.marketplaceUpdate(e.ID)
			}
		}
		return m, nil
	}

	return m, nil
}

// resolveGlobalRoot returns the global ~/.glyph path.
func (m *Model) resolveGlobalRoot() (string, error) {
	ws, err := m.resolver.ResolveGlobal()
	if err != nil {
		return "", err
	}
	return ws.RootPath, nil
}

// resolveProjectRoot returns the project .glyph directory path.
func (m *Model) resolveProjectRoot() string {
	if m.projectConfigPath == "" {
		return ""
	}
	return filepath.Dir(m.projectConfigPath)
}
