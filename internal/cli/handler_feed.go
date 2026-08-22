package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

func HandlerAddFeed(s *state.State, cmd Command) error {
	err := checkArgs(2, []string{"name", "url"}, cmd)
	if err != nil {
		return err
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	currentUser, err := s.Db.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return fmt.Errorf(
			"failed to get current user (%s) from database: %w",
			s.Config.CurrentUserName,
			err,
		)
	}

	feed, err := s.Db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		Name:      name,
		Url:       url,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    currentUser.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Println("Feed added:")
	fmt.Println("	ID:       ", feed.ID)
	fmt.Println("	Name:     ", feed.Name)
	fmt.Println("	URL:      ", feed.Url)
	fmt.Println("	CreateAt: ", feed.CreatedAt)
	fmt.Println("	UpdatedAt:", feed.UpdatedAt)
	fmt.Println("	UserID:   ", feed.UserID)
	return nil
}

func HandlerListFeeds(s *state.State, cmd Command) error {
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
