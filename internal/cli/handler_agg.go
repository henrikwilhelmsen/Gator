package cli

import (
	"context"
	"fmt"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/rss"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
)

func scrapeFeeds(s *state.State) error {
	feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	err = s.Db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return err
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}

	fmt.Print("\n=====================================================================\n\n")
	fmt.Printf("Fetched feed: %s %s\n", feed.Name, feed.Url)
	fmt.Printf("Time:         %s\n\n", time.Now())
	for _, item := range rssFeed.Channel.Item {
		fmt.Println(item.Title)
	}
	fmt.Printf("\nFeed '%s' collected, found %d entries\n", feed.Name, len(rssFeed.Channel.Item))

	return nil
}

// handlerAgg sets up the RSS aggregator
func handlerAgg(s *state.State, cmd Command) error {
	err := checkArgs(1, []string{"time_between_reqs"}, cmd)
	if err != nil {
		return err
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return err
	}

	feeds, err := s.Db.GetFeeds(context.Background())
	if err != nil {
		return err
	}
	fmt.Print("\n=====================================================================\n\n")
	fmt.Print("Feeds to fetch:\n\n")
	for _, f := range feeds {
		fmt.Printf("  %s\n", f.Name)
	}

	fmt.Printf("\nCollecting feeds every %s...\n", cmd.Args[0])
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err := scrapeFeeds(s)
		if err != nil {
			return err
		}
	}
}
