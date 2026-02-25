package main

import (
	"log"

	"github.com/Noudea/glyph/internal/core"
	"github.com/Noudea/glyph/internal/registry"
	"github.com/Noudea/glyph/internal/shell"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rootPath, err := core.DefaultRootPath()
	if err != nil {
		log.Fatal(err)
	}
	state := &core.State{}
	reg := registry.Default()
	model := shell.NewModel(state, rootPath, reg)

	program := tea.NewProgram(model, tea.WithAltScreen())
	reg.SetNotifier(program.Send)
	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
