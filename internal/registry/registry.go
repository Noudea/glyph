package registry

import (
	"sort"
	"strings"
	"sync"

	"github.com/Noudea/glyph/internal/core"
	tasksmodule "github.com/Noudea/glyph/internal/modules/tasks/module"
	tea "github.com/charmbracelet/bubbletea"
)

type Registry struct {
	mu      sync.Mutex
	modules map[string]core.Module
	notify  func(tea.Msg)
}

func New() *Registry {
	return &Registry{modules: make(map[string]core.Module)}
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

func (r *Registry) Register(module core.Module) bool {
	if module == nil {
		return false
	}
	id := module.ID()
	if id == "" {
		return false
	}
	r.mu.Lock()
	if r.modules == nil {
		r.modules = make(map[string]core.Module)
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

func (r *Registry) Modules() []core.Module {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]core.Module, 0, len(r.modules))
	for _, module := range r.modules {
		out = append(out, module)
	}
	return out
}

func (r *Registry) Commands() []core.Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	commands := make([]core.Command, 0, len(r.modules))
	for _, module := range r.modules {
		if module == nil {
			continue
		}
		id := module.ID()
		if id == "" {
			continue
		}
		commands = append(commands, core.Command{
			ID:    id,
			Label: module.Title(),
			Kind:  core.CommandApp,
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
