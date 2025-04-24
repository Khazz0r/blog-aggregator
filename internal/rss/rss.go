package rss

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"log"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// gets xml data from provided url and returns a RSSFeed struct filled in
func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		log.Fatalf("error creating request from context and URL: %v", err)
	}

	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making request: %v", err)
	}
	defer resp.Body.Close()
	

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading body: %v", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(data, &feed)
	if err != nil {
		log.Fatalf("error unmarshalling xml: %v", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Item {
		item := &feed.Channel.Item[i]
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
	}

	return &feed, nil
}