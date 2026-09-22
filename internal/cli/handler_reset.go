// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

// handlerReset resets the state of the program by deleting all records in the database.
// This exists only because we are operating on a toy database to make development easier.
func handlerReset(s *state.State, cmd Command) error {
	err := checkArgs(0, []string{}, cmd)
	if err != nil {
		return err
	}

	// delete all users from the database
	err = s.Db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete users from database: %w", err)
	}

	// delete all feeds from the database
	err = s.Db.DeleteAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to delete feeds from database: %w", err)
	}

	fmt.Println("Successfully reset the program")
	return nil
}
