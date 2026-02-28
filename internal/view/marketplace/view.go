package view

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Noudea/glyph/internal/marketplace"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Entry represents one spellbook row in the marketplace.
type Entry struct {
	ID               string
	Remote           marketplace.Spellbook
	InstalledGlobal  bool
	InstalledProject bool
	HasUpdateGlobal  bool
	HasUpdateProject bool
}

// Installed returns true if installed in any scope.
func (e Entry) Installed() bool {
	return e.InstalledGlobal || e.InstalledProject
}

// HasUpdate returns true if an update is available in any scope.
func (e Entry) HasUpdate() bool {
	return e.HasUpdateGlobal || e.HasUpdateProject
}

// ViewState holds the data the marketplace view needs.
type ViewState struct {
	Loading        bool
	Err            string
	Entries        []Entry
	Cursor         int
	Installing     string
	ConfirmInstall string // non-empty = showing scope prompt for this ID
	HasProject     bool   // whether a project root exists
	Width          int
	Height         int
}

type styles struct {
	title      lipgloss.Style
	count      lipgloss.Style
	muted      lipgloss.Style
	row        lipgloss.Style
	rowActive  lipgloss.Style
	badge      lipgloss.Style
	badgeUpd   lipgloss.Style
	badgeInst  lipgloss.Style
	panel      lipgloss.Style
	accent     lipgloss.Style
	scopeGlob  lipgloss.Style
	scopeProj  lipgloss.Style
	warn       lipgloss.Style
	versionTag lipgloss.Style
	cmdID      lipgloss.Style
	footerKey  lipgloss.Style
	footerDesc lipgloss.Style
	divider    lipgloss.Style
}

func newStyles() styles {
	return styles{
		title:     lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF9F68")),
		count:     lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA3B8")),
		muted:     lipgloss.NewStyle().Foreground(lipgloss.Color("#8A90A6")),
		row:       lipgloss.NewStyle().Foreground(lipgloss.Color("#E7EBF2")),
		rowActive: lipgloss.NewStyle().Foreground(lipgloss.Color("#2F1E0C")).Background(lipgloss.Color("#FFD9A0")).Bold(true),
		badge: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2F1E0C")).
			Background(lipgloss.Color("#A8E6CF")).
			Bold(true).
			Padding(0, 1),
		badgeUpd: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2F1E0C")).
			Background(lipgloss.Color("#FFD3B6")).
			Bold(true).
			Padding(0, 1),
		badgeInst: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#2F1E0C")).
			Background(lipgloss.Color("#A8E6CF")).
			Bold(true).
			Padding(0, 1),
		panel: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#5C6475")).
			Padding(0, 1),
		accent: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9F68")).Bold(true),
		scopeGlob: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A3C9F9")).
			Bold(true),
		scopeProj: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C4B5FD")).
			Bold(true),
		warn:       lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB86C")),
		versionTag: lipgloss.NewStyle().Foreground(lipgloss.Color("#9AA3B8")),
		cmdID:      lipgloss.NewStyle().Foreground(lipgloss.Color("#6B7280")),
		footerKey:  lipgloss.NewStyle().Foreground(lipgloss.Color("#FF9F68")).Bold(true),
		footerDesc: lipgloss.NewStyle().Foreground(lipgloss.Color("#8A90A6")),
		divider:    lipgloss.NewStyle().Foreground(lipgloss.Color("#5C6475")),
	}
}

// Render draws the marketplace panel.
func Render(state ViewState) string {
	panel := renderPanel(state)
	if state.Width > 0 && state.Height > 0 {
		panel = lipgloss.Place(state.Width, state.Height, lipgloss.Center, lipgloss.Center, panel)
	}
	return panel
}

func renderPanel(state ViewState) string {
	s := newStyles()
	panelWidth := resolvePanelWidth(state.Width)
	contentWidth := panelWidth - 4
	if contentWidth < 24 {
		contentWidth = 24
	}

	var b strings.Builder

	// Header: title + count.
	countLabel := fmt.Sprintf("%d spellbooks", len(state.Entries))
	header := joinColumns(
		s.title.Render("✦ Spellbook Marketplace"),
		s.count.Render(countLabel),
		contentWidth,
	)
	b.WriteString(header)
	b.WriteString("\n")
	b.WriteString(s.muted.Render(strings.Repeat("·", contentWidth)))
	b.WriteString("\n")

	if state.Loading {
		b.WriteString(s.muted.Width(contentWidth).Render("Fetching spellbooks..."))
		b.WriteString("\n")
	} else if state.Err != "" {
		b.WriteString(s.warn.Width(contentWidth).Render("⚠ " + state.Err))
		b.WriteString("\n")
	} else if len(state.Entries) == 0 {
		b.WriteString(s.muted.Width(contentWidth).Render("No spellbooks available"))
		b.WriteString("\n")
	} else {
		// Split-pane: left list + divider + right detail.
		leftWidth := contentWidth * 35 / 100
		if leftWidth < 20 {
			leftWidth = 20
		}
		rightWidth := contentWidth - leftWidth - 1 // 1 for divider
		if rightWidth < 20 {
			rightWidth = 20
		}

		leftPanel := renderLeftPanel(state, leftWidth, s)
		rightPanel := renderRightPanel(state, rightWidth, s)

		// Equalize heights.
		leftLines := strings.Count(leftPanel, "\n") + 1
		rightLines := strings.Count(rightPanel, "\n") + 1
		maxLines := leftLines
		if rightLines > maxLines {
			maxLines = rightLines
		}
		for leftLines < maxLines {
			leftPanel += "\n"
			leftLines++
		}
		for rightLines < maxLines {
			rightPanel += "\n"
			rightLines++
		}

		// Build vertical divider.
		dividerCol := strings.Repeat(s.divider.Render("│")+"\n", maxLines)
		dividerCol = strings.TrimRight(dividerCol, "\n")

		body := lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(leftWidth).Render(leftPanel),
			dividerCol,
			lipgloss.NewStyle().Width(rightWidth).PaddingLeft(1).Render(rightPanel),
		)
		b.WriteString(body)
		b.WriteString("\n")
	}

	b.WriteString(renderFooter(state, contentWidth, s))

	return s.panel.Render(strings.TrimRight(b.String(), "\n"))
}

func renderLeftPanel(state ViewState, width int, s styles) string {
	var b strings.Builder
	cursor := clampCursor(state.Cursor, len(state.Entries))
	maxRows := resolveVisibleRows(state.Height, len(state.Entries))
	visible, start, end := entryWindow(state.Entries, cursor, maxRows)

	if start > 0 {
		b.WriteString(s.muted.Render("↑ " + strconv.Itoa(start) + " above"))
		b.WriteString("\n")
	}

	for i, entry := range visible {
		index := start + i
		active := index == cursor
		b.WriteString(renderCompactRow(entry, active, state.Installing, width, s))
		if i < len(visible)-1 {
			b.WriteString("\n")
		}
	}

	if end < len(state.Entries) {
		b.WriteString("\n")
		b.WriteString(s.muted.Render("↓ " + strconv.Itoa(len(state.Entries)-end) + " more"))
	}

	return b.String()
}

func renderCompactRow(e Entry, active bool, installing string, width int, s styles) string {
	prefix := "  "
	if active {
		prefix = "✦ "
	}
	left := prefix + e.Remote.Name

	var badge string
	if installing == e.ID {
		badge = s.badgeUpd.Render("...")
	} else if e.Installed() {
		badge = s.badgeInst.Render("✓")
	}

	row := joinColumns(left, badge, width)

	if active {
		return s.rowActive.Width(width).Render(row)
	}
	return s.row.Width(width).Render(row)
}

func renderRightPanel(state ViewState, width int, s styles) string {
	cursor := clampCursor(state.Cursor, len(state.Entries))
	if cursor >= len(state.Entries) {
		return s.muted.Render("No spellbook selected")
	}
	entry := state.Entries[cursor]
	sb := entry.Remote

	var b strings.Builder

	// Name + scope badges.
	b.WriteString(s.accent.Render(sb.Name))
	if entry.InstalledGlobal {
		b.WriteString("  ")
		b.WriteString(s.badgeInst.Render("✓ global"))
	}
	if entry.InstalledProject {
		b.WriteString("  ")
		b.WriteString(s.badgeInst.Render("✓ project"))
	}
	b.WriteString("\n")

	// Description.
	if sb.Description != "" {
		b.WriteString(s.muted.Render(ansi.Truncate(sb.Description, width, "…")))
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Metadata.
	if sb.Author != "" {
		b.WriteString(s.muted.Render("Author:   ") + sb.Author)
		b.WriteString("\n")
	}
	if sb.Version != "" {
		b.WriteString(s.muted.Render("Version:  ") + sb.Version)
		b.WriteString("\n")
	}
	b.WriteString("\n")

	// Commands.
	if len(sb.Commands) > 0 {
		b.WriteString(s.muted.Render("Commands:"))
		b.WriteString("\n")
		for _, cmd := range sb.Commands {
			label := cmd.Label
			if label == "" {
				label = cmd.ID
			}
			cmdLine := "  " + label + "  " + s.cmdID.Render(cmd.ID)
			b.WriteString(ansi.Truncate(cmdLine, width, "…"))
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func renderFooter(state ViewState, width int, s styles) string {
	// If confirming install scope, show scope prompt.
	if state.ConfirmInstall != "" {
		return s.footerKey.Render("g") + s.footerDesc.Render(" global") +
			s.footerDesc.Render(" · ") +
			s.footerKey.Render("p") + s.footerDesc.Render(" project") +
			s.footerDesc.Render(" · ") +
			s.footerKey.Render("esc") + s.footerDesc.Render(" cancel")
	}

	var parts []string

	// Contextual actions based on selected entry.
	if state.Cursor >= 0 && state.Cursor < len(state.Entries) {
		e := state.Entries[state.Cursor]
		if state.Installing == "" {
			if !e.Installed() {
				parts = append(parts, s.footerKey.Render("i")+s.footerDesc.Render(" install"))
			}
			if e.Installed() {
				parts = append(parts, s.footerKey.Render("u")+s.footerDesc.Render(" uninstall"))
			}
			if e.HasUpdate() {
				parts = append(parts, s.footerKey.Render("U")+s.footerDesc.Render(" update"))
			}
		}
	}

	// Always-present navigation.
	parts = append(parts,
		s.footerKey.Render("esc")+s.footerDesc.Render(" back"),
	)

	footer := strings.Join(parts, s.footerDesc.Render(" · "))
	return ansi.Truncate(footer, width, "…")
}

func resolvePanelWidth(screenWidth int) int {
	const (
		defaultWidth = 120
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

func resolveVisibleRows(height, count int) int {
	if count <= 0 {
		return 0
	}
	rows := 10
	if height > 0 {
		available := height - 8 // header + footer + borders
		if available < 3 {
			available = 3
		}
		if available < rows {
			rows = available
		}
	}
	if rows > count {
		rows = count
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

func entryWindow(entries []Entry, cursor, maxRows int) ([]Entry, int, int) {
	if len(entries) == 0 || maxRows <= 0 {
		return nil, 0, 0
	}
	if len(entries) <= maxRows {
		return entries, 0, len(entries)
	}
	start := cursor - (maxRows / 2)
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > len(entries) {
		end = len(entries)
		start = end - maxRows
		if start < 0 {
			start = 0
		}
	}
	return entries[start:end], start, end
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
