package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

func HandlerFollowing(s *state.State, cmd Command, user database.User) error {
	err := checkArgs(0, []string{}, cmd)
	if err != nil {
		return err
	}

	feeds, err := s.Db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf(
			"failed to get feed follows from database: %w",
			err,
		)
	}

	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}
	return nil
}
