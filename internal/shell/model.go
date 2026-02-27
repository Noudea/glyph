package shell

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Noudea/glyph/internal/core"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Mode int

const (
	ModeSplash Mode = iota
	ModeMain
	ModeLauncher
)

// Model drives the UI.
type Model struct {
	state    *core.State
	resolver core.WorkspaceResolver
	startDir string

	mode Mode

	err string

	launcherInput  textinput.Model
	launcherCursor int

	commandShortcuts map[string][]string
	shortcutCommands map[string]string

	globalConfigPath  string
	projectConfigPath string

	splashFrame int

	width  int
	height int
}

func NewModel(state *core.State, resolver core.WorkspaceResolver) *Model {
	if state == nil {
		state = &core.State{}
	}
	li := textinput.New()
	li.Placeholder = "command"
	li.CharLimit = 64
	li.Width = 24
	li.Prompt = "> "

	startDir := ""
	if rawStart := strings.TrimSpace(resolver.CWD); rawStart != "" {
		startDir = filepath.Clean(rawStart)
	}
	if startDir == "" {
		cwd, err := os.Getwd()
		if err == nil {
			startDir = filepath.Clean(cwd)
		}
	}
	if startDir == "" {
		startDir = "."
	}

	model := &Model{
		state:         state,
		resolver:      resolver,
		startDir:      startDir,
		mode:          ModeSplash,
		launcherInput: li,
	}
	if err := model.reloadConfig(); err != nil {
		model.err = err.Error()
	}
	return model
}

func (m *Model) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 1)
	if m.mode == ModeSplash {
		cmds = append(cmds, splashTickCmd())
	}
	return tea.Batch(cmds...)
}
