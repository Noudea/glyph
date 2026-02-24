package tasks

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type inputMode int

const (
	inputNone inputMode = iota
	inputAdd
	inputEdit
)

type Task struct {
	ID    int
	Title string
	Done  bool
}

type Model struct {
	tasks     []Task
	cursor    int
	input     textinput.Model
	mode      inputMode
	editIndex int
	rootPath  string
}

type ViewModel struct {
	Tasks      []Task
	Cursor     int
	Width      int
	Height     int
	InputView  string
	InputLabel string
	ShowInput  bool
}

func NewModel() Model {
	ti := textinput.New()
	ti.Placeholder = "task title"
	ti.CharLimit = 128
	ti.Width = 32
	ti.Prompt = "> "
	return Model{
		input:     ti,
		editIndex: -1,
	}
}

func (m *Model) SetRoot(path string) error {
	m.rootPath = strings.TrimSpace(path)
	m.cursor = 0
	m.mode = inputNone
	m.editIndex = -1
	m.input.SetValue("")
	m.input.Blur()
	if m.rootPath == "" {
		m.tasks = nil
		return nil
	}
	return m.load()
}

func (m Model) InputFocused() bool {
	return m.input.Focused()
}

func (m Model) Hint() string {
	if m.mode == inputAdd {
		return "tasks: enter add · esc cancel"
	}
	if m.mode == inputEdit {
		return "tasks: enter save · esc cancel"
	}
	return "tasks: a add · e edit · d delete · x toggle"
}

func (m Model) ViewModel(width, height int) ViewModel {
	label := ""
	showInput := m.mode == inputAdd || m.mode == inputEdit
	if m.mode == inputAdd {
		label = "add "
	}
	if m.mode == inputEdit {
		label = "edit "
	}
	return ViewModel{
		Tasks:      m.tasks,
		Cursor:     m.cursor,
		Width:      width,
		Height:     height,
		InputView:  m.input.View(),
		InputLabel: label,
		ShowInput:  showInput,
	}
}

func (m Model) Update(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.mode == inputAdd || m.mode == inputEdit {
		return m.updateInput(msg)
	}
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.tasks)-1 {
			m.cursor++
		}
	case "x", "enter":
		m.toggleTask(m.cursor)
	case "a":
		m.beginAdd()
	case "e":
		m.beginEdit()
	case "d":
		m.deleteTask(m.cursor)
	}
	return m, nil
}

func (m Model) updateInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	switch msg.String() {
	case "esc":
		m.mode = inputNone
		m.editIndex = -1
		m.input.SetValue("")
		m.input.Blur()
	case "enter":
		title := strings.TrimSpace(m.input.Value())
		if title != "" {
			if m.mode == inputAdd {
				m.addTask(title)
			} else if m.mode == inputEdit {
				m.updateTask(title)
			}
		}
		m.mode = inputNone
		m.editIndex = -1
		m.input.SetValue("")
		m.input.Blur()
	}

	return m, cmd
}

func (m *Model) beginAdd() {
	m.mode = inputAdd
	m.input.SetValue("")
	m.input.Focus()
}

func (m *Model) beginEdit() {
	if len(m.tasks) == 0 || m.cursor < 0 || m.cursor >= len(m.tasks) {
		return
	}
	m.mode = inputEdit
	m.editIndex = m.cursor
	m.input.SetValue(m.tasks[m.cursor].Title)
	m.input.Focus()
}

func (m *Model) addTask(title string) {
	nextID := 1
	for _, task := range m.tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}
	m.tasks = append(m.tasks, Task{ID: nextID, Title: title})
	m.cursor = len(m.tasks) - 1
	_ = m.save()
}

func (m *Model) updateTask(title string) {
	if m.editIndex < 0 || m.editIndex >= len(m.tasks) {
		return
	}
	m.tasks[m.editIndex].Title = title
	_ = m.save()
}

func (m *Model) toggleTask(index int) {
	if index < 0 || index >= len(m.tasks) {
		return
	}
	m.tasks[index].Done = !m.tasks[index].Done
	_ = m.save()
}

func (m *Model) deleteTask(index int) {
	if index < 0 || index >= len(m.tasks) {
		return
	}
	m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
	if m.cursor >= len(m.tasks) && len(m.tasks) > 0 {
		m.cursor = len(m.tasks) - 1
	}
	if len(m.tasks) == 0 {
		m.cursor = 0
	}
	_ = m.save()
}

func (m *Model) load() error {
	path := tasksFilePath(m.rootPath)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			m.tasks = nil
			return nil
		}
		return err
	}
	var items []Task
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	m.tasks = items
	if m.cursor >= len(m.tasks) {
		m.cursor = len(m.tasks) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	return nil
}

func (m *Model) save() error {
	if m.rootPath == "" {
		return nil
	}
	path := tasksFilePath(m.rootPath)
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(m.tasks, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func tasksFilePath(rootPath string) string {
	return filepath.Join(rootPath, "tasks", "tasks.json")
}
