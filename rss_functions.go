package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	RSSstruct := &RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	req.Header.Set(req.UserAgent(), "gator")
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail when making request: %w", err)
	}

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail getting response: %w", err)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("fail when reading response: %w", err)
	}
	xml.Unmarshal(data, RSSstruct)
	return RSSstruct, nil
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

func aggregate() error {
	resp, err := http.Get("https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("fail when fetching data: %w", err)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading data: %w", err)
	}
	fmt.Println(data)
	return nil
}
