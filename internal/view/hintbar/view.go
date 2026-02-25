package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type ViewState struct {
	Width int
	Text  string
}

func Render(state ViewState) string {
	if state.Text == "" {
		return ""
	}

	text := strings.TrimSpace(state.Text)
	content := text
	if state.Width > 0 {
		maxContent := state.Width - 4
		if maxContent > 0 && lipgloss.Width(text) > maxContent {
			text = ansi.Truncate(text, maxContent, "…")
		}
		if maxContent > 0 {
			content = lipgloss.NewStyle().Width(maxContent).Align(lipgloss.Center).Render(text)
		}
	}

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7A7A7A"))
	if state.Width > 0 {
		style = style.Width(state.Width)
	}

	return style.Render(content)
}

func Height(width int) int {
	return 1
}
