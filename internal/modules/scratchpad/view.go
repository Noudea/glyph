package scratchpad

import (
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

func Render(state ViewModel) string {
	lines := make([]string, 0, 2)
	if state.Error != "" {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B"))
		lines = append(lines, errStyle.Render("error: "+state.Error))
	}
	if state.Editing {
		lines = append(lines, state.EditorView)
	} else {
		lines = append(lines, state.PreviewView)
	}
	content := strings.TrimRight(strings.Join(lines, "\n"), "\n")
	style := lipgloss.NewStyle()
	if state.Width > 0 {
		style = style.Width(state.Width)
	}
	if state.Height > 0 {
		style = style.Height(state.Height)
	}
	return style.Render(content)
}

func renderMarkdown(content string, width int) string {
	if strings.TrimSpace(content) == "" {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A")).Render("empty scratchpad")
	}

	wrap := width
	if wrap < 20 {
		wrap = 80
	}

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(wrap),
	)
	if err != nil {
		return content
	}
	out, err := renderer.Render(content)
	if err != nil {
		return content
	}
	return strings.TrimRight(out, "\n")
}
