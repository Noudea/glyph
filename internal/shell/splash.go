package shell

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	splashFrameInterval = 90 * time.Millisecond
	splashFrameCount    = 16
)

type splashTickMsg struct{}

func splashTickCmd() tea.Cmd {
	return tea.Tick(splashFrameInterval, func(time.Time) tea.Msg {
		return splashTickMsg{}
	})
}

func (m *Model) finishSplash() {
	m.splashFrame = 0
	m.openLauncher()
}
