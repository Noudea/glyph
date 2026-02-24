package app

import tea "github.com/charmbracelet/bubbletea"

// Module defines the contract for a built-in or external app.
// Keep this interface stable to support plugins later.
type Module interface {
	ID() string
	Title() string
	Init(ctx AppContext) tea.Cmd
	Update(ctx AppContext, msg tea.Msg) (Module, tea.Cmd)
	View(width, height int) string
	Hint() string
	InputFocused() bool
}

// AppContext carries shared runtime context for modules.
type AppContext struct {
	RootPath string
}
