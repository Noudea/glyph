package shell

import (
	"sort"

	"github.com/Noudea/glyph/internal/app"
	"github.com/Noudea/glyph/internal/modules/registry"
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
	state *app.State
	root  string

	mode Mode

	registry *registry.Registry
	modules  map[string]app.Module
	err      string

	launcherInput  textinput.Model
	launcherCursor int

	width  int
	height int
}

func NewModel(state *app.State, rootPath string, registry *registry.Registry) *Model {
	if state == nil {
		state = &app.State{}
	}
	li := textinput.New()
	li.Placeholder = "command"
	li.CharLimit = 64
	li.Width = 24
	li.Prompt = "> "

	moduleIndex := make(map[string]app.Module)
	if registry != nil {
		moduleIndex, _ = indexModules(registry.Modules())
	}
	if state.Commands == nil && registry != nil {
		state.Commands = registry.Commands()
	}
	if len(state.OpenApps) == 0 && len(state.Commands) > 0 {
		state.OpenApps = []string{state.Commands[0].ID}
	}
	if state.ActiveApp == "" && len(state.OpenApps) > 0 {
		state.ActiveApp = state.OpenApps[0]
	}

	return &Model{
		state:         state,
		root:          rootPath,
		mode:          ModeMain,
		registry:      registry,
		modules:       moduleIndex,
		launcherInput: li,
	}
}

func (m *Model) Init() tea.Cmd {
	return moduleInitCmd(m.modules, m.context())
}

func moduleInitCmd(modules map[string]app.Module, ctx app.AppContext) tea.Cmd {
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

func indexModules(modules []app.Module) (map[string]app.Module, []app.Command) {
	moduleIndex := make(map[string]app.Module, len(modules))
	commands := make([]app.Command, 0, len(modules))
	for _, module := range modules {
		if module == nil {
			continue
		}
		id := module.ID()
		if id == "" {
			continue
		}
		moduleIndex[id] = module
		commands = append(commands, app.Command{
			ID:    id,
			Label: module.Title(),
			Kind:  app.CommandApp,
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
	if len(m.state.OpenApps) == 0 && len(commands) > 0 {
		m.state.OpenApps = []string{commands[0].ID}
	}
	if m.state.ActiveApp == "" && len(m.state.OpenApps) > 0 {
		m.state.ActiveApp = m.state.OpenApps[0]
	}
}
