package main

import (
	"log"
	"os"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	resolver := core.NewWorkspaceResolver(cwd)
	state := &core.State{}
	model := shell.NewModel(state, resolver)

	program := tea.NewProgram(model, tea.WithAltScreen())
	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
