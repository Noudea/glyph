package view

import (
	"strconv"
	"strings"

	"github.com/Noudea/glyph/internal/core"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type ViewState struct {
	InputView string
	Commands  []core.Command
	Cursor    int
	Width     int
	Height    int
}

type paletteStyles struct {
	title        lipgloss.Style
	count        lipgloss.Style
	muted        lipgloss.Style
	searchBox    lipgloss.Style
	row          lipgloss.Style
	rowActive    lipgloss.Style
	shortcutChip lipgloss.Style
	activeChip   lipgloss.Style
	panel        lipgloss.Style
}

func Render(state ViewState) string {
	panel := renderPanel(state)

	if state.Width > 0 && state.Height > 0 {
		panel = lipgloss.Place(state.Width, state.Height, lipgloss.Center, lipgloss.Center, panel)
	}
	return panel
}

func renderPanel(state ViewState) string {
	styles := newPaletteStyles()
	panelWidth := resolvePanelWidth(state.Width)
	contentWidth := panelWidth - 4
	if contentWidth < 24 {
		contentWidth = 24
	}

	cursor := clampCursor(state.Cursor, len(state.Commands))
	maxRows := resolveVisibleRows(state.Height, len(state.Commands))
	visible, start, end := commandWindow(state.Commands, cursor, maxRows)

	var b strings.Builder

	header := joinColumns(
		styles.title.Render("✦ Spellbook"),
		styles.count.Render(strconv.Itoa(len(state.Commands))+" spells"),
		contentWidth,
	)
	b.WriteString(header)
	b.WriteString("\n")
	searchWidth := contentWidth - 4
	if searchWidth < 1 {
		searchWidth = 1
	}
	b.WriteString(styles.searchBox.Width(searchWidth).Render(state.InputView))
	b.WriteString("\n")
	b.WriteString(styles.muted.Render(strings.Repeat("·", contentWidth)))
	b.WriteString("\n")

	if len(state.Commands) == 0 {
		b.WriteString(styles.muted.Width(contentWidth).Render("No matching commands"))
	} else {
		if start > 0 {
			topMore := styles.muted.Render("↑ " + strconv.Itoa(start) + " above")
			b.WriteString(styles.muted.Width(contentWidth).Render(topMore))
			b.WriteString("\n")
		}

		for i, cmd := range visible {
			index := start + i
			active := index == cursor
			b.WriteString(renderCommandRow(cmd, active, contentWidth, styles))
			if i < len(visible)-1 || end < len(state.Commands) {
				b.WriteString("\n")
			}
		}

		if end < len(state.Commands) {
			bottomMore := styles.muted.Render("↓ " + strconv.Itoa(len(state.Commands)-end) + " more")
			if len(visible) > 0 {
				b.WriteString("\n")
			}
			b.WriteString(styles.muted.Width(contentWidth).Render(bottomMore))
		}
	}

	b.WriteString("\n")
	b.WriteString(styles.muted.Width(contentWidth).Render("enter cast · esc close · ↑/↓ navigate"))

	return styles.panel.Render(strings.TrimRight(b.String(), "\n"))
}

func newPaletteStyles() paletteStyles {
	return paletteStyles{
		title: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF9F68")),
		count: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#9AA3B8")),
		muted: lipgloss.NewStyle().Foreground(lipgloss.Color("#8A90A6")),
		searchBox: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E7EBF2")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5C6475")).
			Padding(0, 1),
		row: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E7EBF2")),
		rowActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2F1E0C")).
			Background(lipgloss.Color("#FFD9A0")).
			Bold(true),
		shortcutChip: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DFEE")).
			Background(lipgloss.Color("#252A35")).
			Padding(0, 1),
		activeChip: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2F1E0C")).
			Background(lipgloss.Color("#FFCF92")).
			Bold(true).
			Padding(0, 1),
		panel: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#5C6475")).
			Padding(0, 1),
	}
}

func resolvePanelWidth(screenWidth int) int {
	const (
		defaultWidth = 88
		minWidth     = 40
	)
	if screenWidth <= 0 {
		return defaultWidth
	}
	max := screenWidth - 4
	if max < minWidth {
		return max
	}
	if defaultWidth > max {
		return max
	}
	return defaultWidth
}

func resolveVisibleRows(height, commandCount int) int {
	if commandCount <= 0 {
		return 0
	}
	rows := 10
	if height > 0 {
		available := height - 9
		if available < 3 {
			available = 3
		}
		if available < rows {
			rows = available
		}
	}
	if rows > commandCount {
		rows = commandCount
	}
	if rows < 1 {
		rows = 1
	}
	return rows
}

func clampCursor(cursor, size int) int {
	if size == 0 {
		return 0
	}
	if cursor < 0 {
		return 0
	}
	if cursor >= size {
		return size - 1
	}
	return cursor
}

func commandWindow(commands []core.Command, cursor, maxRows int) ([]core.Command, int, int) {
	if len(commands) == 0 || maxRows <= 0 {
		return nil, 0, 0
	}
	if len(commands) <= maxRows {
		return commands, 0, len(commands)
	}
	start := cursor - (maxRows / 2)
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > len(commands) {
		end = len(commands)
		start = end - maxRows
		if start < 0 {
			start = 0
		}
	}
	return commands[start:end], start, end
}

func renderCommandRow(cmd core.Command, active bool, width int, styles paletteStyles) string {
	left := "  " + cmd.Label
	if active {
		left = "✦ " + cmd.Label
	}

	meta := make([]string, 0, 2)
	if cmd.Shortcut != "" {
		if active {
			meta = append(meta, styles.activeChip.Render(cmd.Shortcut))
		} else {
			meta = append(meta, styles.shortcutChip.Render(cmd.Shortcut))
		}
	}
	right := strings.Join(meta, " ")

	row := joinColumns(left, right, width)
	if active {
		return styles.rowActive.Width(width).Render(row)
	}
	return styles.row.Width(width).Render(row)
}

func joinColumns(left, right string, width int) string {
	if width < 1 {
		return left
	}
	if right == "" {
		return ansi.Truncate(left, width, "…")
	}

	rightWidth := lipgloss.Width(right)
	maxLeft := width - rightWidth - 1
	if maxLeft < 1 {
		return ansi.Truncate(left+" "+right, width, "…")
	}

	left = ansi.Truncate(left, maxLeft, "…")
	leftWidth := lipgloss.Width(left)
	space := width - leftWidth - rightWidth
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right
}
