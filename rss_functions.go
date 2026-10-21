package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/muhammad-romli/ArachNet/internal/database"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssStruct := &RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("fail when making request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail getting response: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("fail when reading response: %w", err)
	}
	if err := xml.Unmarshal(data, rssStruct); err != nil {
		return nil, fmt.Errorf("fail unmarshaling xml data: %w", err)
	}
	return rssStruct, nil
}

func processXML(feed *RSSFeed) {
	items := feed.Channel.Item
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := 0; i < len(items); i++ {
		items[i].Title = html.UnescapeString(items[i].Title)
		items[i].Description = html.UnescapeString(items[i].Description)
	}
}

func scrapeFeeds(s *state) error {
	dbFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("fail getting new feed to fetch: %w", err)
	}
	s.db.MarkFeedFetched(context.Background(), dbFeed.ID)
	rssFeed, err := fetchFeed(context.Background(), dbFeed.Url)
	if err != nil {
		return err
	}
	processXML(rssFeed)
	for _, item := range rssFeed.Channel.Item {
		publishedAt := sql.NullTime{} // Valid: false, stored as NULL
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		}
		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: publishedAt,
			FeedID:      dbFeed.ID,
		}
		s.db.CreatePost(context.Background(), params)
	}
	return nil
}

func stripHTML(input string) string { // Because some sites like hacker news doesn't have description
	re := regexp.MustCompile("<[^>]*>")
	return re.ReplaceAllString(input, "")
}
