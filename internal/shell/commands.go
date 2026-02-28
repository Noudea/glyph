package shell

import (
	"sort"
	"strings"

	"github.com/Noudea/glyph/internal/core"
)

const (
	commandSourceGlobal    = "global"
	commandSourceProject   = "project"
	commandSourceManaged   = "managed"
	commandSourceSpellbook = "spellbook"
)

func (m Model) launcherCommands() []core.Command {
	out := make([]core.Command, 0, len(m.state.Commands)+1)

	// Add synthetic marketplace command.
	out = append(out, core.Command{
		ID:      commandMarketplaceOpen,
		Label:   "Spellbook Marketplace",
		Kind:    core.CommandAction,
		Group:   "system",
		Source:  commandSourceManaged,
		Managed: true,
	})

	if m.state == nil || len(m.state.Commands) == 0 {
		return out
	}

	for _, command := range m.state.Commands {
		item := command
		item.Shortcut = m.primaryShortcut(item.ID, item.Shortcut)
		out = append(out, item)
	}

	sort.SliceStable(out, func(i, j int) bool {
		left := strings.ToLower(strings.TrimSpace(out[i].Label))
		right := strings.ToLower(strings.TrimSpace(out[j].Label))
		if left == right {
			return out[i].ID < out[j].ID
		}
		return left < right
	})

	return out
}
