package view

import (
	"strings"

	"github.com/Noudea/glyph/internal/modules/tasks"
	"github.com/charmbracelet/lipgloss"
)

func Render(state tasks.ViewModel) string {
	var b strings.Builder
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A"))
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E2E2")).Bold(true)
	done := lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A")).Strikethrough(true)

	if len(state.Tasks) == 0 {
		b.WriteString(muted.Render("no tasks yet"))
	} else {
		for i, task := range state.Tasks {
			indicator := "  "
			style := lipgloss.NewStyle().Foreground(lipgloss.Color("#E2E2E2"))
			label := "• " + task.Title
			if task.Done {
				label = "✓ " + task.Title
				style = done
			}
			if i == state.Cursor {
				indicator = "› "
				style = active
				if task.Done {
					style = active.Foreground(lipgloss.Color("#7A7A7A")).Strikethrough(true)
				}
			}
			b.WriteString(indicator)
			b.WriteString(style.Render(label))
			b.WriteString("\n")
		}
	}

	if state.ShowInput {
		if b.Len() > 0 {
			b.WriteString("\n")
		}
		label := state.InputLabel
		if label == "" {
			label = "add "
		}
		b.WriteString(muted.Render(label))
		b.WriteString(state.InputView)
	}

	content := strings.TrimRight(b.String(), "\n")
	if state.Width > 0 || state.Height > 0 {
		style := lipgloss.NewStyle()
		if state.Width > 0 {
			style = style.Width(state.Width)
		}
		if state.Height > 0 {
			style = style.Height(state.Height)
		}
		return style.Render(content)
	}
	return content
}
