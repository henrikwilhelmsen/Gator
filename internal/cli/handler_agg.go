package cli

import (
	"context"
	"fmt"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/rss"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

// HandlerAgg sets up the RSS aggregator
func HandlerAgg(s *state.State, cmd Command) error {
	err := checkArgs(0, []string{}, cmd)
	if err != nil {
		return err
	}

	// fetch a single feed for testing and lesson submission
	url := "https://www.wagslane.dev/index.xml"
	feed, err := rss.FetchFeed(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	fmt.Println(feed)
	return nil
}
