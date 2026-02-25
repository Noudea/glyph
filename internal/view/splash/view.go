package view

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ViewState struct {
	Frame  int
	Width  int
	Height int
}

var sparkFrames = []string{
	"✦   ·   ✧",
	"· ✶   ✦ ·",
	"✧   ✦   ·",
	"· ✦   ✶ ·",
}

func Render(state ViewState) string {
	lines := make([]string, 0, 4)
	styles := splashStyles()

	lines = append(lines, styles.spark.Render(sparkFrames[state.Frame%len(sparkFrames)]))
	lines = append(lines, "")
	lines = append(lines, styles.title.Render("glyph"))
	lines = append(lines, styles.subtitle.Render("opening the spellbook..."))

	content := strings.Join(lines, "\n")
	if state.Width > 0 && state.Height > 0 {
		return lipgloss.Place(state.Width, state.Height, lipgloss.Center, lipgloss.Center, content)
	}
	return content
}

type styles struct {
	spark    lipgloss.Style
	title    lipgloss.Style
	subtitle lipgloss.Style
}

func splashStyles() styles {
	return styles{
		spark: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFCF92")).
			Bold(true),
		title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9F68")).
			Bold(true),
		subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
	}
}
