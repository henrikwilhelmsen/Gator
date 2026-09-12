package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

// Wrapper function to convert a handler that takes a user to a regular handler.
func middlewareLoggedIn(
	handler func(s *state.State, cmd Command, user database.User) error,
) func(*state.State, Command) error {
	// Return a function that does not need a user argument.
	return func(s *state.State, cmd Command) error {
		// Get the current user and return an error if it fails.
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return fmt.Errorf(
				"failed to get current user (%s) from database: %w",
				s.Config.CurrentUserName,
				err,
			)
		}

		// Return the wrapped function with the user argument supplied
		return handler(s, cmd, user)
	}
}
