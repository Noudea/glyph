package core

// CommandKind describes the command type used by the launcher.
type CommandKind string

const (
	CommandApp    CommandKind = "app"
	CommandAction CommandKind = "action"
	CommandExec   CommandKind = "exec"
)

// Command describes a runnable command in the palette.
type Command struct {
	ID       string
	Label    string
	Kind     CommandKind
	Group    string
	Shortcut string
	Run      string
	Source   string
	Managed  bool
	ToolID   string
}

// State holds shared app state across UI.
type State struct {
	ActiveApp string
	OpenApps  []string
	Commands  []Command
	Workspace Workspace
}
