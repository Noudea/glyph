package shell

import (
	"sort"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/registry"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeMain Mode = iota
	ModeLauncher
)

// Model drives the UI.
type Model struct {
	state *core.State
	root  string

	workspace core.Workspace
	resolver  core.WorkspaceResolver

	mode Mode

	registry *registry.Registry
	modules  map[string]core.Module
	err      string

	launcherInput  textinput.Model
	launcherCursor int

	commandShortcuts map[string][]string
	shortcutCommands map[string]string

	width  int
	height int
}

func NewModel(state *core.State, workspace core.Workspace, resolver core.WorkspaceResolver, registry *registry.Registry) *Model {
	if state == nil {
		state = &core.State{}
	}
	li := textinput.New()
	li.Placeholder = "command"
	li.CharLimit = 64
	li.Width = 24
	li.Prompt = "> "

	moduleIndex := make(map[string]core.Module)
	commands := make([]core.Command, 0)
	if registry != nil {
		moduleIndex, commands = indexModules(registry.Modules())
	}
	if registry != nil {
		state.Commands = commands
	}
	if workspace.RootPath == "" {
		workspace = state.Workspace
	}
	state.Workspace = workspace

	model := &Model{
		state:         state,
		root:          workspace.RootPath,
		workspace:     workspace,
		resolver:      resolver,
		mode:          ModeLauncher,
		registry:      registry,
		modules:       moduleIndex,
		launcherInput: li,
	}
	if err := model.loadGlobalShortcuts(); err != nil {
		model.err = err.Error()
	}
	model.openLauncher()
	return model
}

func (m *Model) Init() tea.Cmd {
	return moduleInitCmd(m.modules, m.context())
}

func moduleInitCmd(modules map[string]core.Module, ctx core.CoreContext) tea.Cmd {
	if len(modules) == 0 {
		return nil
	}
	cmds := make([]tea.Cmd, 0, len(modules))
	for _, module := range modules {
		if module == nil {
			continue
		}
		cmds = append(cmds, module.Init(ctx))
	}
	return tea.Batch(cmds...)
}

func indexModules(modules []core.Module) (map[string]core.Module, []core.Command) {
	moduleIndex := make(map[string]core.Module, len(modules))
	commands := make([]core.Command, 0, len(modules))
	for _, module := range modules {
		if module == nil {
			continue
		}
		id := module.ID()
		if id == "" {
			continue
		}
		moduleIndex[id] = module
		commands = append(commands, core.Command{
			ID:    id,
			Label: module.Title(),
			Kind:  core.CommandApp,
			Group: "apps",
		})
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Label < commands[j].Label
	})
	return moduleIndex, commands
}

func (m *Model) refreshRegistry() {
	if m.registry == nil {
		return
	}
	moduleIndex, commands := indexModules(m.registry.Modules())
	m.modules = moduleIndex
	if m.state == nil {
		return
	}
	m.state.Commands = commands
	validOpen := make([]string, 0, len(m.state.OpenApps))
	for _, id := range m.state.OpenApps {
		if _, ok := moduleIndex[id]; ok {
			validOpen = append(validOpen, id)
		}
	}
	m.state.OpenApps = validOpen
	if m.state.ActiveApp != "" {
		if _, ok := moduleIndex[m.state.ActiveApp]; !ok {
			m.state.ActiveApp = ""
		}
	}
	if m.state.ActiveApp == "" && len(m.state.OpenApps) > 0 {
		m.state.ActiveApp = m.state.OpenApps[0]
	}
}

func (m *Model) setWorkspace(workspace core.Workspace) tea.Cmd {
	if workspace.RootPath == "" {
		return nil
	}
	if m.workspace.Kind == workspace.Kind &&
		m.workspace.RootPath == workspace.RootPath &&
		m.workspace.ProjectPath == workspace.ProjectPath {
		return nil
	}

	m.workspace = workspace
	m.root = workspace.RootPath
	if m.state != nil {
		m.state.Workspace = workspace
	}
	if err := m.loadGlobalShortcuts(); err != nil {
		m.err = err.Error()
	} else {
		m.err = ""
	}
	return moduleInitCmd(m.modules, m.context())
}
