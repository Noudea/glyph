package core

// CommandKind describes the command type used by the launcher.
type CommandKind string

const (
	CommandApp    CommandKind = "app"
	CommandAction CommandKind = "action"
)

// Command describes a runnable command in the palette.
type Command struct {
	ID       string
	Label    string
	Kind     CommandKind
	Group    string
	Shortcut string
}

// State holds shared app state across UI.
type State struct {
	ActiveApp string
	OpenApps  []string
	Commands  []Command
	Workspace Workspace
}
