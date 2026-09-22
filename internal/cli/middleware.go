// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

// middlewareLoggedIn is a wrapper function to convert a handler that requires a user to a regular handler.
func middlewareLoggedIn(
	handler func(s *state.State, cmd Command, user database.User) error,
) func(*state.State, Command) error {
	return func(s *state.State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}
