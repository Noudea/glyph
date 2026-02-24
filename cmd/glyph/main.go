package main

import (
	"log"

	"github.com/Noudea/glyph/internal/app"
	"github.com/Noudea/glyph/internal/modules/registry"
	"github.com/Noudea/glyph/internal/modules/shell"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	rootPath, err := app.DefaultRootPath()
	if err != nil {
		log.Fatal(err)
	}
	state := &app.State{}
	reg := registry.Default()
	model := shell.NewModel(state, rootPath, reg)

	program := tea.NewProgram(model, tea.WithAltScreen())
	reg.SetNotifier(program.Send)
	if err := program.Start(); err != nil {
		log.Fatal(err)
	}
}
