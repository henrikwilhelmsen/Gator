package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/config"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"github.com/google/uuid"
)

type State struct {
	Config *config.Config
	Db     database.Querier
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	// Mapping of command names to handler functions
	CommandsToHandlers map[string]func(*State, Command) error
}

// Run executes the command with the given state if a handler can be located.
func (c *Commands) Run(s *State, cmd Command) error {
	handler, ok := c.CommandsToHandlers[cmd.Name]
	if !ok {
		return fmt.Errorf("No handler for command '%s' registered", cmd.Name)
	}
	return handler(s, cmd)
}

// Register registers the given name to the given handler function.
func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.CommandsToHandlers[name] = f
}

// LoginHandler sets the current user to the given commands single argument. The
// cmd must have exactly one argument, the user name. Anything else will return error.
func HandlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(
			"Command '%s' expects exactly 1 argument (username), got %d",
			cmd.Name,
			len(cmd.Args))
	}

	// Will error and return if the user has not been added to the database
	_, err := s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return err
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	fmt.Printf("Current user set to %s\n", cmd.Args[0])
	return nil
}

// HandlerRegister registers the given username in the database.
func HandlerRegister(s *State, cmd Command) error {
	// Check that we only have one argument
	if len(cmd.Args) != 1 {
		return fmt.Errorf(
			"Command '%s' expects exactly 1 argument (username), got %d",
			cmd.Name,
			len(cmd.Args))
	}

	// Create the user in the database
	usr, err := s.Db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        uuid.New(),
			Name:      cmd.Args[0],
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	)
	if err != nil {
		return err
	}

	// Update the config
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}

	// Print result and return nil
	fmt.Printf("Registered user: %v\n", usr)
	return nil
}

// HandlerReset resets the state of the program by deleting all records in the database.
// This exists only because we are operating on a toy database to make development easier.
func HandlerReset(s *State, cmd Command) error {
	err := s.Db.DeleteAll(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Successfully reset database")
	return nil
}

// SetupRegisterCommands sets up and returns a Commands struct, with all supported
// commands registered. The current commands are: login
func SetupRegisterCommands() *Commands {
	commands := Commands{
		CommandsToHandlers: make(map[string]func(*State, Command) error),
	}
	commands.Register("login", HandlerLogin)
	commands.Register("register", HandlerRegister)
	commands.Register("reset", HandlerReset)
	return &commands
}
