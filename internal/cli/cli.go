package cli

import (
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
)

type State struct {
	Config *config.Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	// Mapping of command names to handler functions
	CommandsToHandlers map[string]func(*State, Command) error
}

// Run the command with the given state if a handler can be located.
func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.CommandsToHandlers[cmd.Name]
	if !ok {
		return fmt.Errorf("No handler for command '%s' registered", cmd.Name)
	}
	return handler(s, cmd)
}

// Register a new command handler
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandsToHandlers[name] = f
}

// Login handler, set the current user to the given argument
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(
			"Command '%s' expects exactly 1 argument (username), got %d",
			cmd.Name,
			len(cmd.Args))
	}

	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %s\n", cmd.Args[0])
	return nil
}
