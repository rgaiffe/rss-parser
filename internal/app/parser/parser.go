package parser

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v9"
	"github.com/mmcdole/gofeed/rss"
	"github.com/rgaiffe/rss-parser/internal/pkg/store"
	"github.com/spf13/viper"
)

// getRssFeed is a function to get the rss feed from config file
func getRssFeed() map[string]string {
	return viper.GetViper().GetStringMapString("rssFeed")
}

// rssParser parse the rss feed
func getFeedsFromRssParser(body io.ReadCloser) (*rss.Feed, error) {
	fp := rss.Parser{}

	feeds, err := fp.Parse(body)
	if err != nil {
		return nil, err
	}
	log.Printf("Feed: %s\n", feeds.Link)

	return feeds, nil
}

// checkIfFeedExist is a function to check if the feed exist in redis
func checkIfFeedExist(client *store.Client, feed *rss.Item) bool {
	if err := client.Get(context.Background(), feed.Link).Err(); errors.Is(err, redis.Nil) {
		return false
	}
	return true
}

// putFeedInStore is a function to put the feed in redis
func putFeedInStore(client *store.Client, feed *rss.Item) error {
	if err := client.Set(context.Background(), feed.Link, time.Now().UTC(), 0).Err(); err != nil {
		return err
	}
	return nil
}

// callRssFeed is a function to make http call to the rss feed
func rssParser(feedUri string) error {
	// Make http call to the rss feed
	ret, err := http.Get(feedUri)
	if err != nil {
		return err
	}
	defer ret.Body.Close()

	if ret.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Error: %d", ret.StatusCode))
	}
	log.Printf("Success: %d\n", ret.StatusCode)

	// Parse the rss feed from the http call
	feeds, err := getFeedsFromRssParser(ret.Body)
	if err != nil {
		return err
	}

	// Create a new redis client
	client, err := store.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new discord client
	discordClient, _ := discordgo.New("")

	for _, feed := range feeds.Items {
		if checkIfFeedExist(client, feed) {
			continue
		}

		webhookParams := &discordgo.WebhookParams{
			Content: feed.Link,
		}

		// Send the message to the discord channel
		_, err = discordClient.WebhookExecute(viper.GetString("discord.webhook.id"), viper.GetString("discord.webhook.token"), true, webhookParams)
		if err != nil {
			return err
		}

		// Set the feed link in redis
		if putFeedInStore(client, feed) != nil {
			return err
		}
	}

	return nil
}

// StartParser is the main function for start the parser
func StartParser() error {
	for k, v := range getRssFeed() {
		fmt.Printf("Starting call: %s: %s\n", k, v)
		if err := rssParser(v); err != nil {
			return err
		}
	}
	return nil
}
