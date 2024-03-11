package scraper

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/DavAnders/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func FetchRSSFeed(url string) (*RSS, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rss RSS
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	return &rss, nil
}

func fetchAndProcessFeed(ctx context.Context, feed database.Feed, queries *database.Queries) {

	feedData, err := FetchRSSFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed %s: %v", feed.Url, err)
		return
	}

	for _, item := range feedData.Channel.Items {
		var parsedTime time.Time
		var parseErr error
		for _, layout := range []string{time.RFC1123, time.RFC1123Z, "2006-01-02T15:04:05Z07:00"} {
			parsedTime, parseErr = time.Parse(layout, item.PubDate)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			log.Printf("Error parsing published date for item: '%s': %v", item.Title, parseErr)
			continue
		}
		err = queries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("Failed to save post '%s': %v", item.Title, err)
			continue
		}
	}
	// marks as fetched whether saved to db or not. could add more granular control over this later
	err = queries.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched %d: %v", feed.ID, err)
	} else {
		log.Printf("Feed with ID %s marked at fetched.", feed.ID)
	}
}

func Worker(ctx context.Context, interval time.Duration, batchSize int32, queries *database.Queries) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Worked stopped")
			return
		case <-time.After(interval):
			log.Println("Fetching feeds...")

			// Fetch the next batch of feeds from the database
			feeds, err := queries.GetNextFeedsToFetch(ctx, batchSize)
			if err != nil {
				log.Printf("Error fetching feeds: %v", err)
				continue
			}
			if len(feeds) == 0 {
				log.Println("No feeds to fetch")
				continue
			}
			log.Printf("Fetched %d feeds to process", len(feeds))

			var wg sync.WaitGroup
			for _, feed := range feeds {
				wg.Add(1)
				go func(feed database.Feed) {
					defer wg.Done()
					fetchAndProcessFeed(ctx, feed, queries)
				}(feed)
			}
			wg.Wait()
			log.Println("Completed fetching feeds.")

		}
	}
}
