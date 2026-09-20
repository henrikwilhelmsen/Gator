package cli

import (
	"context"
	"fmt"
	"strconv"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

// handlerBrowse prints the stored posts for the current user, with an optional `limit`
// arg used to limit the amount of posts to show - defaults to 2.
func handlerBrowse(s *state.State, cmd Command, user database.User) error {
	var limit int32 = 2
	if len(cmd.Args) == 1 {
		limitArg, err := strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = int32(limitArg)
	}

	posts, err := s.Db.GetPostsForUser(
		context.Background(),
		database.GetPostsForUserParams{UserID: user.ID, Limit: limit},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to get posts from database: %w",
			err,
		)
	}

	for i, post := range posts {
		fmt.Print("\n=====================================================================\n\n")
		fmt.Printf("%s from %s\n", post.PublishedAt.Time.Format("Mon Jan 2"), post.FeedName)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Description: %v\n", post.Description.String)
		fmt.Printf("URL: %s\n", post.Url)
		fmt.Println()
		fmt.Println(i + 1)
	}
	return nil
}
