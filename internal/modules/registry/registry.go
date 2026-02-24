package registry

import (
	"sort"
	"strings"
	"sync"

	"github.com/Noudea/glyph/internal/app"
	tasksmodule "github.com/Noudea/glyph/internal/modules/tasks/module"
	tea "github.com/charmbracelet/bubbletea"
)

type Registry struct {
	mu      sync.Mutex
	modules map[string]app.Module
	notify  func(tea.Msg)
}

func New() *Registry {
	return &Registry{modules: make(map[string]app.Module)}
}

func Default() *Registry {
	r := New()
	_ = r.Register(tasksmodule.NewModule())
	return r
}

// ChangedMsg is emitted when the registry contents change.
type ChangedMsg struct{}

func (r *Registry) SetNotifier(notify func(tea.Msg)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notify = notify
}

func (r *Registry) Register(module app.Module) bool {
	if module == nil {
		return false
	}
	id := module.ID()
	if id == "" {
		return false
	}
	r.mu.Lock()
	if r.modules == nil {
		r.modules = make(map[string]app.Module)
	}
	if _, exists := r.modules[id]; exists {
		r.mu.Unlock()
		return false
	}
	r.modules[id] = module
	notify := r.notify
	r.mu.Unlock()
	if notify != nil {
		notify(ChangedMsg{})
	}
	return true
}

func (r *Registry) Modules() []app.Module {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]app.Module, 0, len(r.modules))
	for _, module := range r.modules {
		out = append(out, module)
	}
	return out
}

func (r *Registry) Commands() []app.Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	commands := make([]app.Command, 0, len(r.modules))
	for _, module := range r.modules {
		if module == nil {
			continue
		}
		id := module.ID()
		if id == "" {
			continue
		}
		commands = append(commands, app.Command{
			ID:    id,
			Label: module.Title(),
			Kind:  app.CommandApp,
			Group: "apps",
		})
	}
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Label < commands[j].Label
	})
	return commands
}

func (r *Registry) Unregister(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" {
		return false
	}
	r.mu.Lock()
	if r.modules == nil {
		r.mu.Unlock()
		return false
	}
	if _, exists := r.modules[id]; !exists {
		r.mu.Unlock()
		return false
	}
	delete(r.modules, id)
	notify := r.notify
	r.mu.Unlock()
	if notify != nil {
		notify(ChangedMsg{})
	}
	return true
}
