// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

func handlerListFeeds(s *state.State, cmd Command) error {
	err := checkArgs(0, []string{}, cmd)
	if err != nil {
		return err
	}
	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get feeds from database: %w", err)
	}
	for _, feed := range feeds {
		user, err := s.Db.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("failed to get username from database: %w", err)
		}

		fmt.Println("Feed:")
		fmt.Println("	Name:    ", feed.Name)
		fmt.Println("	URL:     ", feed.Url)
		fmt.Println("	Added By:", user.Name)
	}
	return nil
}
