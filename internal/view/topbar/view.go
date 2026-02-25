package view

import (
	"strconv"
	"strings"

	"github.com/Noudea/glyph/internal/core"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type ViewState struct {
	Width     int
	Title     string
	Tabs      []core.Command
	ActiveApp string
}

type topbarStyles struct {
	line      lipgloss.Style
	brand     lipgloss.Style
	separator lipgloss.Style
	tabActive lipgloss.Style
	tabIdle   lipgloss.Style
	overflow  lipgloss.Style
	emptyTabs lipgloss.Style
}

func Render(state ViewState) string {
	if state.Title == "" {
		return ""
	}

	const horizontalPadding = 1
	styles := newTopbarStyles()
	left := styles.brand.Render("✦ glyph · " + displayTitle(state.Title))

	const gap = "  "
	if state.Width > 0 {
		availableTabs := state.Width - (horizontalPadding * 2) - lipgloss.Width(left) - lipgloss.Width(gap)
		if availableTabs < 1 {
			line := renderLine(left, state.Width, styles.line, horizontalPadding)
			spacer := renderLine("", state.Width, styles.line, horizontalPadding)
			return spacer + "\n" + line + "\n" + spacer
		}
		tabs := renderTabsContent(availableTabs, state.Tabs, state.ActiveApp, styles)
		content := lipgloss.JoinHorizontal(lipgloss.Left, left, styles.separator.Render(gap), tabs)
		line := renderLine(content, state.Width, styles.line, horizontalPadding)
		spacer := renderLine("", state.Width, styles.line, horizontalPadding)
		return spacer + "\n" + line + "\n" + spacer
	}

	tabs := renderTabsContent(0, state.Tabs, state.ActiveApp, styles)
	content := lipgloss.JoinHorizontal(lipgloss.Left, left, styles.separator.Render(gap), tabs)
	return "\n" + styles.line.Render(" "+content+" ") + "\n"
}

func newTopbarStyles() topbarStyles {
	return topbarStyles{
		line: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8ECF6")),
		brand: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#2B1A0F")).
			Background(lipgloss.Color("#FF9F68")).
			Padding(0, 1),
		separator: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
		tabActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#2B1A0F")).
			Background(lipgloss.Color("#FFCF92")).
			Padding(0, 1),
		tabIdle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#D8DFEE")).
			Background(lipgloss.Color("#252A35")).
			Padding(0, 1),
		overflow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
		emptyTabs: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
	}
}

func renderTabsContent(width int, tabs []core.Command, activeID string, styles topbarStyles) string {
	if len(tabs) == 0 {
		return styles.emptyTabs.Render("no open apps")
	}

	chips := make([]string, 0, len(tabs))
	for _, tab := range tabs {
		style := styles.tabIdle
		if tab.ID == activeID {
			style = styles.tabActive
		}
		chips = append(chips, style.Render(tab.Label))
	}

	if width <= 0 {
		return strings.Join(chips, " ")
	}

	out := make([]string, 0, len(chips))
	used := 0
	for i, chip := range chips {
		chipWidth := lipgloss.Width(chip)
		sepWidth := 0
		if len(out) > 0 {
			sepWidth = 1
		}

		remaining := len(chips) - i - 1
		overflow := ""
		overflowWidth := 0
		if remaining > 0 {
			overflow = styles.overflow.Render("+" + strconv.Itoa(remaining))
			overflowWidth = 1 + lipgloss.Width(overflow)
		}

		if used+sepWidth+chipWidth+overflowWidth > width {
			if len(out) == 0 {
				return ansi.Truncate(chip, width, "…")
			}
			out = append(out, overflow)
			break
		}

		out = append(out, chip)
		used += sepWidth + chipWidth
	}

	return strings.Join(out, " ")
}

func displayTitle(title string) string {
	title = strings.TrimSpace(title)
	title = strings.TrimPrefix(title, "glyph · ")
	if title == "" {
		return "global"
	}
	return title
}

func renderLine(content string, width int, style lipgloss.Style, horizontalPadding int) string {
	if width <= 0 {
		return style.Render(strings.Repeat(" ", horizontalPadding) + content + strings.Repeat(" ", horizontalPadding))
	}
	if horizontalPadding < 0 {
		horizontalPadding = 0
	}
	innerWidth := width - (horizontalPadding * 2)
	if innerWidth < 1 {
		innerWidth = 1
	}
	if lipgloss.Width(content) > innerWidth {
		content = ansi.Truncate(content, innerWidth, "…")
	}
	padding := innerWidth - lipgloss.Width(content)
	if padding < 0 {
		padding = 0
	}
	line := strings.Repeat(" ", horizontalPadding) + content + strings.Repeat(" ", padding) + strings.Repeat(" ", horizontalPadding)
	return style.Render(line)
}

func Height(width int) int {
	return 3
}
