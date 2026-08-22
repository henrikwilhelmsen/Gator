package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

// HandlerRegister registers the given username in the database.
func HandlerRegister(s *state.State, cmd Command) error {
	err := checkArgs(1, []string{"username"}, cmd)
	if err != nil {
		return err
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
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update the config
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	// Print result and return nil
	fmt.Printf("Registered user: %v\n", usr)
	return nil
}

// LoginHandler sets the current user to the given commands single argument. The
// cmd must have exactly one argument, the user name. Anything else will return error.
func HandlerLogin(s *state.State, cmd Command) error {
	err := checkArgs(1, []string{"username"}, cmd)
	if err != nil {
		return err
	}

	// Will error and return if the user has not been added to the database
	_, err = s.Db.GetUser(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("user not found in database, unable to log in: %w", err)
	}

	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("failed to set user: %w", err)
	}

	fmt.Printf("Current user set to: %s\n", cmd.Args[0])
	return nil
}

// HandlerListUsers print all of the users to the console
func HandlerListUsers(s *state.State, cmd Command) error {
	err := checkArgs(0, []string{}, cmd)
	if err != nil {
		return err
	}
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get users from database: %w", err)
	}
	for _, usr := range users {
		if usr.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", usr.Name)
		} else {
			fmt.Printf("* %s\n", usr.Name)
		}
	}
	return nil
}
