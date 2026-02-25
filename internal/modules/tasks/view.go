package tasks

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

type taskStyles struct {
	title      lipgloss.Style
	stats      lipgloss.Style
	divider    lipgloss.Style
	row        lipgloss.Style
	rowDone    lipgloss.Style
	rowTodo    lipgloss.Style
	rowActive  lipgloss.Style
	iconDone   lipgloss.Style
	iconTodo   lipgloss.Style
	empty      lipgloss.Style
	overflow   lipgloss.Style
	inputLabel lipgloss.Style
	inputBox   lipgloss.Style
}

func Render(state ViewModel) string {
	styles := newTaskStyles()
	width := state.Width

	total := len(state.Tasks)
	doneCount := countDone(state.Tasks)
	openCount := total - doneCount

	header := joinColumns(
		styles.title.Render("✦ Tasks"),
		styles.stats.Render(fmt.Sprintf("%d open · %d done", openCount, doneCount)),
		width,
	)
	divider := styles.divider.Render(renderDivider(width))

	availableRows := visibleTaskRows(state.Height, state.ShowInput)
	visibleTasks, start, end := taskWindow(state.Tasks, state.Cursor, availableRows)

	var b strings.Builder
	b.WriteString(header)
	b.WriteString("\n")
	b.WriteString(divider)
	b.WriteString("\n")

	if total == 0 {
		b.WriteString(styles.empty.Render("No tasks yet. Press a to add one."))
	} else {
		if start > 0 {
			b.WriteString(styles.overflow.Render("↑ " + fmt.Sprintf("%d above", start)))
			b.WriteString("\n")
		}

		for i, task := range visibleTasks {
			taskIndex := start + i
			active := taskIndex == state.Cursor
			b.WriteString(renderTaskRow(task, active, width, styles))
			if i < len(visibleTasks)-1 || end < total {
				b.WriteString("\n")
			}
		}

		if end < total {
			if len(visibleTasks) > 0 {
				b.WriteString("\n")
			}
			b.WriteString(styles.overflow.Render("↓ " + fmt.Sprintf("%d more", total-end)))
		}
	}

	if state.ShowInput {
		b.WriteString("\n\n")
		label := strings.TrimSpace(state.InputLabel)
		if label == "" {
			label = "add"
		}
		b.WriteString(styles.inputLabel.Render(strings.ToUpper(label)))
		b.WriteString("\n")
		b.WriteString(styles.inputBox.Render(state.InputView))
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

func newTaskStyles() taskStyles {
	return taskStyles{
		title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF9F68")),
		stats: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
		divider: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5C6475")),
		row: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8ECF6")),
		rowDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A1A9BC")),
		rowTodo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#E8ECF6")),
		rowActive: lipgloss.NewStyle().
			Bold(true),
		iconDone: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF9F68")).
			Bold(true),
		iconTodo: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5C6475")).
			Bold(true),
		empty: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
		overflow: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A90A6")),
		inputLabel: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF9F68")),
		inputBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#5C6475")).
			Padding(0, 1),
	}
}

func renderTaskRow(task Task, active bool, width int, styles taskStyles) string {
	marker := "◇"
	if task.Done {
		marker = "◆"
	}

	prefix := "  "
	if active {
		prefix = "✶ "
	}

	iconStyle := styles.iconTodo
	textStyle := styles.rowTodo
	if task.Done {
		iconStyle = styles.iconDone
		textStyle = styles.rowDone
	}
	content := prefix + iconStyle.Render(marker) + " " + textStyle.Render(task.Title)

	if width > 0 {
		content = ansi.Truncate(content, width, "…")
	}

	style := styles.row
	if active {
		style = styles.rowActive
	}
	if width > 0 {
		style = style.Width(width)
	}
	return style.Render(content)
}

func visibleTaskRows(height int, showInput bool) int {
	// Header consumes 2 rows. Input area consumes 3 rows when visible.
	rows := 10
	if height > 0 {
		rows = height - 2
		if showInput {
			rows -= 3
		}
	}
	if rows < 1 {
		rows = 1
	}
	return rows
}

func taskWindow(tasks []Task, cursor, rows int) ([]Task, int, int) {
	if len(tasks) == 0 || rows <= 0 {
		return nil, 0, 0
	}
	if rows >= len(tasks) {
		return tasks, 0, len(tasks)
	}

	if cursor < 0 {
		cursor = 0
	}
	if cursor >= len(tasks) {
		cursor = len(tasks) - 1
	}

	start := cursor - (rows / 2)
	if start < 0 {
		start = 0
	}
	end := start + rows
	if end > len(tasks) {
		end = len(tasks)
		start = end - rows
		if start < 0 {
			start = 0
		}
	}
	return tasks[start:end], start, end
}

func countDone(tasks []Task) int {
	count := 0
	for _, t := range tasks {
		if t.Done {
			count++
		}
	}
	return count
}

func renderDivider(width int) string {
	if width <= 0 {
		return "··········"
	}
	return strings.Repeat("·", width)
}

func joinColumns(left, right string, width int) string {
	if width <= 0 || right == "" {
		return left
	}
	rightWidth := lipgloss.Width(right)
	maxLeft := width - rightWidth - 1
	if maxLeft < 1 {
		return ansi.Truncate(left, width, "…")
	}
	left = ansi.Truncate(left, maxLeft, "…")
	space := width - lipgloss.Width(left) - rightWidth
	if space < 1 {
		space = 1
	}
	return left + strings.Repeat(" ", space) + right
}
