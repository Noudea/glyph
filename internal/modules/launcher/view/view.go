package view

import (
	"strings"

	"github.com/Noudea/glyph/internal/app"
	"github.com/charmbracelet/lipgloss"
)

type ViewState struct {
	InputView string
	Commands  []app.Command
	Cursor    int
	Width     int
	Height    int
}

func Render(state ViewState) string {
	panel := renderPanel(state)

	if state.Width > 0 && state.Height > 0 {
		panel = lipgloss.Place(state.Width, state.Height, lipgloss.Center, lipgloss.Center, panel)
	}
	return panel
}

func renderPanel(state ViewState) string {
	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Render("command palette")
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A"))
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD166")).Bold(true)

	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(state.InputView)
	b.WriteString("\n")
	b.WriteString(muted.Render("────────────────────────────"))
	b.WriteString("\n")

	for i, cmd := range state.Commands {
		indicator := "  "
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E2E2"))
		if i == state.Cursor {
			indicator = "> "
			style = active
		}
		b.WriteString(indicator)
		b.WriteString(style.Render(cmd.Label))
		b.WriteString("\n")
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		Padding(0, 1)

	return box.Render(b.String())
}
