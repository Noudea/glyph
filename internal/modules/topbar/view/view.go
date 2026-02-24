package view

import (
	"strconv"
	"strings"

	"github.com/Noudea/glyph/internal/app"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type ViewState struct {
	Width     int
	Title     string
	Tabs      []app.Command
	ActiveApp string
}

func Render(state ViewState) string {
	if state.Title == "" {
		return ""
	}

	workspace := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E2E2E2")).
		Bold(true).
		Render("⬡ " + state.Title)

	parts := buildTopBarParts(state.Width, workspace, state.Tabs, state.ActiveApp)
	return renderTopBarLine(state.Width, parts)
}

func buildTopBarParts(width int, workspace string, tabs []app.Command, activeID string) []string {
	active := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD166")).Bold(true)
	inactive := lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A"))

	parts := []string{workspace}
	if len(tabs) == 0 {
		return parts
	}

	if width <= 0 {
		for _, tab := range tabs {
			style := inactive
			if tab.ID == activeID {
				style = active
			}
			parts = append(parts, style.Render(tab.Label))
		}
		return parts
	}

	maxContent := width - 4
	if maxContent < 1 {
		return parts
	}

	sepWidth := lipgloss.Width(" │ ")
	contentWidth := lipgloss.Width(workspace)
	for i, tab := range tabs {
		style := inactive
		if tab.ID == activeID {
			style = active
		}
		part := style.Render(tab.Label)
		partWidth := lipgloss.Width(part)

		remaining := len(tabs) - i - 1
		overflowWidth := 0
		if remaining > 0 {
			overflowWidth = sepWidth + lipgloss.Width(renderOverflow(remaining))
		}
		needed := sepWidth + partWidth + overflowWidth
		if contentWidth+needed > maxContent {
			parts = append(parts, renderOverflow(remaining+1))
			return parts
		}
		parts = append(parts, part)
		contentWidth += sepWidth + partWidth
	}
	return parts
}

func renderOverflow(count int) string {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("#7A7A7A"))
	return style.Render("›› " + strconv.Itoa(count) + " more")
}

func renderTopBarLine(width int, parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	sep := " │ "
	content := strings.Join(parts, sep)
	if width <= 0 {
		return content
	}

	if width < 4 {
		return content
	}

	maxContent := width - 4
	if lipgloss.Width(content) > maxContent {
		content = ansi.Truncate(content, maxContent, "…")
	}
	padding := maxContent - lipgloss.Width(content)
	if padding < 0 {
		padding = 0
	}
	line := "│ " + content + strings.Repeat(" ", padding) + " │"
	top := "┌" + strings.Repeat("─", width-2) + "┐"
	bot := "└" + strings.Repeat("─", width-2) + "┘"
	return strings.Join([]string{top, line, bot}, "\n")
}

func Height(width int) int {
	if width <= 0 {
		return 1
	}
	return 3
}
