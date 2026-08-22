package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/rss"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

// LoginHandler sets the current user to the given commands single argument. The
// cmd must have exactly one argument, the user name. Anything else will return error.
func HandlerLogin(s *state.State, cmd Command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf(
			"command '%s' expects exactly 1 argument (username), got %d",
			cmd.Name,
			len(cmd.Args),
		)
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
func HandlerRegister(s *state.State, cmd Command) error {
	// Check that we only have one argument
	if len(cmd.Args) != 1 {
		return fmt.Errorf(
			"command '%s' expects exactly 1 argument (username), got %d",
			cmd.Name,
			len(cmd.Args),
		)
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
func HandlerReset(s *state.State, cmd Command) error {
	err := s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return err
	}
	err = s.Db.DeleteAllFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Println("Successfully reset database")
	return nil
}

// HandlerUsers print all of the users to the console
func HandlerUsers(s *state.State, cmd Command) error {
	users, err := s.Db.GetUsers(context.Background())
	if err != nil {
		return err
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

// HandlerAgg sets up the RSS aggregator
func HandlerAgg(s *state.State, cmd Command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf(
			"command '%s' expects no arguments, got %d",
			cmd.Name,
			len(cmd.Args),
		)
	}

	// fetch a single feed for testing and lesson submission
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return err
	}

	fmt.Println(feed)
	return nil
}

func HandlerAddFeed(s *state.State, cmd Command) error {
	// get current user and connect to the feed
	return nil
}
