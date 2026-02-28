package marketplace

// Spellbook describes a community spellbook manifest.
type Spellbook struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Version     string    `json:"version"`
	Commands    []Command `json:"commands"`
}

// Command describes a single command within a spellbook.
// Same shape as config's commandConfig.
type Command struct {
	ID      string `json:"id"`
	Label   string `json:"label"`
	Run     string `json:"run,omitempty"`
	Script  string `json:"script,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}
