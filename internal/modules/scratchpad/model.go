package scratchpad

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	modePreview mode = iota
	modeEdit
)

type Model struct {
	rootPath      string
	content       string
	mode          mode
	input         textarea.Model
	err           string
	scroll        int
	previewHeight int
}

type ViewModel struct {
	Width       int
	Height      int
	Editing     bool
	EditorView  string
	PreviewView string
	Error       string
}

func NewModel() Model {
	input := textarea.New()
	input.Placeholder = ""
	input.ShowLineNumbers = false
	input.CharLimit = 0
	input.Blur()
	return Model{
		mode:  modePreview,
		input: input,
	}
}

func (m *Model) SetRoot(path string) error {
	m.rootPath = strings.TrimSpace(path)
	m.mode = modePreview
	m.err = ""
	m.scroll = 0
	m.previewHeight = 0
	m.input.SetValue("")
	m.input.Blur()
	if m.rootPath == "" {
		m.content = ""
		return nil
	}
	return m.load()
}

func (m Model) Hint() string {
	if m.mode == modeEdit {
		return "scratchpad: ctrl+s save · esc cancel"
	}
	return "scratchpad: ↑/↓ scroll · e edit"
}

func (m *Model) ViewModel(width, height int) ViewModel {
	contentHeight := height
	if contentHeight < 1 {
		contentHeight = 1
	}
	if m.err != "" {
		contentHeight--
		if contentHeight < 1 {
			contentHeight = 1
		}
	}
	m.previewHeight = contentHeight

	state := ViewModel{
		Width:   width,
		Height:  height,
		Editing: m.mode == modeEdit,
		Error:   m.err,
	}
	if m.mode == modeEdit {
		m.input.SetHeight(contentHeight)
		if width > 0 {
			m.input.SetWidth(width)
		}
		state.EditorView = m.input.View()
		return state
	}
	preview := renderMarkdown(m.content, width)
	state.PreviewView, m.scroll = previewWindow(preview, m.scroll, contentHeight)
	return state
}

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.mode == modeEdit {
		return m.updateEdit(msg)
	}
	switch msg.String() {
	case "e":
		m.beginEdit()
	case "down", "j":
		m.scroll++
	case "up", "k":
		if m.scroll > 0 {
			m.scroll--
		}
	case "pgdown":
		m.scroll += m.pageStep()
	case "pgup":
		step := m.pageStep()
		m.scroll -= step
		if m.scroll < 0 {
			m.scroll = 0
		}
	case "home", "g":
		m.scroll = 0
	case "end", "G":
		// Clamped during rendering based on content size.
		m.scroll = int(^uint(0) >> 1)
	}
	return m, nil
}

func (m Model) updateEdit(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = modePreview
		m.input.SetValue(m.content)
		m.input.Blur()
		return m, nil
	case "ctrl+s":
		m.content = m.input.Value()
		if err := m.save(); err != nil {
			m.err = err.Error()
			return m, nil
		}
		m.err = ""
		m.mode = modePreview
		m.input.Blur()
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *Model) beginEdit() {
	m.mode = modeEdit
	m.input.SetValue(m.content)
	m.input.Focus()
}

func (m *Model) load() error {
	path := scratchpadFilePath(m.rootPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				m.content = ""
				return err
			}
			if err := os.WriteFile(path, nil, 0o644); err != nil {
				m.content = ""
				return err
			}
			m.content = ""
			m.input.SetValue("")
			return nil
		}
		m.content = ""
		return err
	}
	m.content = string(data)
	m.input.SetValue(m.content)
	m.scroll = 0
	return nil
}

func (m *Model) save() error {
	if m.rootPath == "" {
		return nil
	}
	path := scratchpadFilePath(m.rootPath)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(m.content), 0o644)
}

func scratchpadFilePath(rootPath string) string {
	return filepath.Join(rootPath, "scratchpad", "scratchpad.md")
}

func (m Model) pageStep() int {
	step := m.previewHeight - 1
	if step < 1 {
		step = 1
	}
	return step
}

func previewWindow(content string, scroll, height int) (string, int) {
	if height < 1 {
		height = 1
	}
	lines := strings.Split(content, "\n")
	maxScroll := len(lines) - height
	if maxScroll < 0 {
		maxScroll = 0
	}
	if scroll < 0 {
		scroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	end := scroll + height
	if end > len(lines) {
		end = len(lines)
	}
	return strings.Join(lines[scroll:end], "\n"), scroll
}
