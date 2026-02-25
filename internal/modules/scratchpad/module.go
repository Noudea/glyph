package scratchpad

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Noudea/glyph/internal/core"
)

type Module struct {
	model    Model
	rootPath string
}

func NewModule() core.Module {
	return &Module{model: NewModel()}
}

func (m *Module) ID() string {
	return "scratchpad"
}

func (m *Module) Title() string {
	return "Scratchpad"
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
	return Render(m.model.ViewModel(width, height))
}

func (m *Module) Hint() string {
	return m.model.Hint()
}

func (m *Module) ensureContext(ctx core.CoreContext) {
	if m.rootPath == ctx.RootPath {
		return
	}
	m.rootPath = ctx.RootPath
	_ = m.model.SetRoot(ctx.RootPath)
}
