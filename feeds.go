package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Khazz0r/blog-aggregator/internal/database"
	"github.com/Khazz0r/blog-aggregator/internal/rss"
)

// grab all items from the next feed to fetch, mark that feed, and print out all of the titles to the console
func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Fatalf("error with trying to retrieve next feed: %v", err)
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{
			Time: time.Now().UTC(),
			Valid: true,
		},
		UpdatedAt: time.Now().UTC(),
		ID: feed.ID,
	})
	if err != nil {
		log.Fatalf("error with trying to mark feed's fetched time: %v", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Fatalf("error fetching items from feed: %v", err)
	}

	for _, item := range rssFeed.Channel.Item {
		fmt.Printf("* title:          %s\n", item.Title)
		fmt.Printf("* link:           %s\n", item.Link)
		fmt.Printf("* description:    %s\n", item.Description)
		fmt.Println("=======================================================================")
	}

	return nil
}
