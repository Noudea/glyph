package core

import tea "github.com/charmbracelet/bubbletea"

// Module defines the contract for a built-in or external module.
// Keep this interface stable to support plugins later.
type Module interface {
	ID() string
	Title() string
	Init(ctx CoreContext) tea.Cmd
	Update(ctx CoreContext, msg tea.Msg) (Module, tea.Cmd)
	View(width, height int) string
	Hint() string
}

// CoreContext carries shared runtime context for modules.
type CoreContext struct {
	RootPath string
}
