// Copyright 2026 Henrik Wilhelmsen. All rights reserved.
// SPDX-License-Identifier: MPL-2.0

package rss

import (
	"context"
	"encoding/xml"
	"errors"
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

func FetchFeed(ctx context.Context, feedURL string) (feed *RSSFeed, err error) {
	client := http.Client{}
	feed = &RSSFeed{}

	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return feed, err
	}
	req.Header.Set("user-agent", "gator")

	// run the request
	res, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	// use errors.Join so that we catch any errors in the Close method.
	defer func() {
		err = errors.Join(err, res.Body.Close())
	}()

	// read the data from the request body
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return feed, err
	}

	// unmarshal the data
	err = xml.Unmarshal(data, feed)
	if err != nil {
		return feed, err
	}

	// un-escape all of the content in the returned data
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}
	return feed, err
}
