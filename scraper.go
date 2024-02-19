package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/nickmarks00/go-rss-aggregator/internal/database"
)

func startScraping(
	db *database.Queries, // connection to database
	concurrency int, // how many Go routines we want running
	timeBetweenRequests time.Duration,
) {
	log.Printf("Scraping on %v go routines every %s duration...", concurrency, timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C { // empty initialisation means that the loop runs immediately on first hit
		feeds, err := db.GetNextFeedsToFetch(
			context.Background(),
			int32(concurrency),
		)
		if err != nil {
			log.Println("error fetching feeds: ", err)
			continue
		}

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1) // add +1 to counter of spawned routines

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait() // wait for counter of spawned routines to zero
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done() // each .Done() call decrements the counter that we had INCREMENTED with Add()

	_, err := db.MarkFeedAsFecthed(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("error fetching feed: ", err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Println("Error converting published date to string: ", err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to create post: ", err)
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
