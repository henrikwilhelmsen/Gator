package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

func handlerFollow(s *state.State, cmd Command, user database.User) error {
	err := checkArgs(1, []string{"url"}, cmd)
	if err != nil {
		return err
	}

	// Get the feed
	url := cmd.Args[0]
	feed, err := s.Db.GetFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf(
			"failed to get feed from database, make sure it has been added: %w",
			err,
		)
	}

	// Create a new feed follow entry
	feedFollowRow, err := s.Db.CreateFeedFollow(
		context.Background(),
		database.CreateFeedFollowParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create feed follow entry in database: %w", err)
	}

	// Print the result
	fmt.Printf(
		"Created feed follow entry for feed '%s' and user '%s'",
		feedFollowRow.FeedName,
		feedFollowRow.UserName,
	)

	return nil
}
