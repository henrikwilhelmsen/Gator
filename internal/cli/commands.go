package cli

import (
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	// Mapping of command names to handler functions
	CommandsToHandlers map[string]func(*state.State, Command) error
}

// Run executes the command with the given state if a handler can be located.
func (c *Commands) Run(s *state.State, cmd Command) error {
	handler, ok := c.CommandsToHandlers[cmd.Name]
	if !ok {
		return fmt.Errorf("no handler for command '%s' registered", cmd.Name)
	}
	return handler(s, cmd)
}

// Register registers the given name to the given handler function.
func (c *Commands) Register(name string, f func(*state.State, Command) error) {
	c.CommandsToHandlers[name] = f
}
