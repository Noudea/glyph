package main

import (
	"log"
	"os"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/registry"
	"github.com/Noudea/glyph/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	resolver := core.NewWorkspaceResolver(cwd)
	workspace, err := resolver.ResolveInitial()
	if err != nil {
		log.Fatal(err)
	}
	state := &core.State{Workspace: workspace}
	reg := registry.Default()
	model := shell.NewModel(state, workspace, resolver, reg)

	program := tea.NewProgram(model, tea.WithAltScreen())
	reg.SetNotifier(program.Send)
	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
