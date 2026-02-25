package module

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/modules/tasks"
	tasksview "github.com/Noudea/glyph/internal/modules/tasks/view"
)

type Module struct {
	model    tasks.Model
	rootPath string
}

func NewModule() core.Module {
	return &Module{model: tasks.NewModel()}
}

func (m *Module) ID() string {
	return "tasks"
}

func (m *Module) Title() string {
	return "Tasks"
}

func (m *Module) Init(ctx core.CoreContext) tea.Cmd {
	m.ensureContext(ctx)
	return nil
}

func (m *Module) Update(ctx core.CoreContext, msg tea.Msg) (core.Module, tea.Cmd) {
	m.ensureContext(ctx)
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	var cmd tea.Cmd
	m.model, cmd = m.model.Update(keyMsg)
	return m, cmd
}

func (m *Module) View(width, height int) string {
	return tasksview.Render(m.model.ViewModel(width, height))
}

func (m *Module) Hint() string {
	return m.model.Hint()
}

func (m *Module) InputFocused() bool {
	return m.model.InputFocused()
}

func (m *Module) ensureContext(ctx core.CoreContext) {
	if m.rootPath == ctx.RootPath {
		return
	}
	m.rootPath = ctx.RootPath
	_ = m.model.SetRoot(ctx.RootPath)
}
