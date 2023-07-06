package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/yeyo27/rssAggregator/internal/database"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequests time.Duration) {
	log.Printf(
		"Scraping on %v go routines every %s duration",
		concurrency,
		timeBetweenRequests)
	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("error fetching feeds:", err)
			continue
		}

		waitGroup := sync.WaitGroup{}
		for _, feed := range feeds {
			waitGroup.Add(1)
			go scrapeFeed(db, &waitGroup, feed)
		}
		waitGroup.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched:", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed:", err)
	}
	for _, item:= range rssFeed.Channel.Item{
		log.Println("Found post:", item.Title, "on feed", feed.Name)
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}