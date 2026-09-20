package cli

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/database"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/rss"
	"git.hwanimation.tech/henrikwilhelmsen/gator/internal/state"
	"github.com/google/uuid"
)

// formatPostTimeString takes a string representing a time and tries to parse it with
// time.Parse, using various common layouts. If parsing a layout succeeds, a valid
// sql.NullTime is returned, otherwise return an invalid sql.NullTime.
func formatPostTimeString(timeString string) sql.NullTime {
	layouts := [6]string{
		time.UnixDate,
		time.RubyDate,
		time.RFC1123Z,
		time.RFC3339,
		time.DateTime,
		time.DateOnly,
	}

	for _, layout := range layouts {
		postTime, err := time.Parse(layout, timeString)
		if err == nil {
			return sql.NullTime{Time: postTime, Valid: true}
		}
	}
	return sql.NullTime{}
}

// scrapeFeeds checks the database for the next feed to fetch, then fetches that feed
// and saves any new posts to the database.
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
		// Convert item data to sql nullable types
		publishedAt := formatPostTimeString(item.PubDate)
		description := sql.NullString{String: item.Description, Valid: true}

		// Create the post
		post, err := s.Db.CreatePost(
			context.Background(),
			database.CreatePostParams{
				ID:          uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Title:       item.Title,
				Url:         item.Link,
				Description: description,
				PublishedAt: publishedAt,
				FeedID:      feed.ID,
			},
		)

		if err != nil {
			// If the error is a duplicate key, ignore the error, otherwise log it
			if !strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("Error when adding post: %s", err)
			}
			continue
		}

		fmt.Printf("Added post '%s' from '%s'\n", post.Title, post.Url)
	}

	fmt.Printf("\nFeed '%s' collected, found %d entries\n", feed.Name, len(rssFeed.Channel.Item))
	return nil
}

// handlerAgg sets up and runs the RSS aggregator
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
