package shell

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type commandFinishedMsg struct {
	CommandID string
	Err       error
}

func (m *Model) executeCommand(commandID string) tea.Cmd {
	switch commandID {
	case commandLauncherOpen:
		m.openLauncher()
		return nil
	case commandMarketplaceOpen:
		return m.openMarketplace()
	default:
		command, ok := m.findCommandByID(commandID)
		if !ok {
			m.err = "command not found: " + commandID
			return nil
		}
		run := strings.TrimSpace(command.Run)
		if run == "" {
			m.err = "command has no run value: " + commandID
			return nil
		}

		process := shellExecCommand(run, m.startDir)
		return tea.ExecProcess(process, func(err error) tea.Msg {
			return commandFinishedMsg{CommandID: commandID, Err: err}
		})
	}
}

func (m *Model) handleCommandFinished(msg commandFinishedMsg) {
	if msg.Err == nil {
		m.err = ""
		m.openLauncher()
		return
	}

	command, ok := m.findCommandByID(msg.CommandID)
	label := msg.CommandID
	if ok {
		label = command.Label
	}
	m.err = fmt.Sprintf("%s failed: %s", label, formatCommandExecError(msg.Err))
	m.openLauncher()
}

func (m *Model) openLauncher() {
	m.mode = ModeLauncher
	m.launcherInput.SetValue("")
	m.launcherInput.Focus()
	m.clampLauncherCursor()
}

func shellExecCommand(run string, cwd string) *exec.Cmd {
	var command *exec.Cmd
	if runtime.GOOS == "windows" {
		command = exec.Command("cmd", "/C", wrapWindowsQuickPauseCommand(run))
	} else {
		command = exec.Command("sh", "-lc", wrapPosixQuickPauseCommand(run))
	}
	command.Dir = cwd
	return command
}

func wrapPosixQuickPauseCommand(run string) string {
	return "clear; __glyph_start=$(date +%s); " +
		run +
		"; __glyph_status=$?; __glyph_end=$(date +%s); " +
		"if [ $((__glyph_end-__glyph_start)) -lt 2 ]; then " +
		"printf '\\n[glyph] Press Enter to return...'; IFS= read -r _; clear; " +
		"fi; exit $__glyph_status"
}

func wrapWindowsQuickPauseCommand(run string) string {
	return "cls & " + run + " & set __glyph_status=%errorlevel% & echo. & echo [glyph] Press Enter to return... & pause >nul & cls & exit /b %__glyph_status%"
}

func formatCommandExecError(err error) string {
	if err == nil {
		return ""
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		code := exitErr.ExitCode()
		if code >= 0 {
			return fmt.Sprintf("exit code %d", code)
		}
	}
	return err.Error()
}
