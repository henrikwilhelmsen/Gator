// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

func handlerUnfollow(s *state.State, cmd Command, user database.User) error {
	err := checkArgs(1, []string{"url"}, cmd)
	if err != nil {
		return err
	}

	// Get the feed
	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf(
			"failed to get feed from database: %w",
			err,
		)
	}

	// Delete the feed follow
	err = s.Db.DeleteFeedFollow(
		context.Background(),
		database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feed.ID},
	)
	if err != nil {
		return fmt.Errorf("failed to create feed follow entry in database: %w", err)
	}

	fmt.Printf("Unfollowed feed: %s\n", feed.Name)
	return nil
}
