package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/gotify/plugin-api"
	"github.com/mmcdole/gofeed"
)

type Storage struct {
	LastPublished time.Time `json:"last_published"`
}

type Config struct {
	RefreshInterval int      `yaml:"refresh_interval"`
	FeedURLs        []string `yaml:"feed_urls"`
}

func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		ModulePath:  "github.com/solarkennedy/gotify-rss",
		Version:     "0.0.1",
		Author:      "solarkennedy",
		Website:     "https://github.com/solarkennedy/gotify-rss",
		Description: "Poll RSS Feeds for Notifications",
		License:     "GPLv3+",
		Name:        "rss",
	}
}

type RssPlugin struct {
	msgHandler     plugin.MessageHandler
	storageHandler plugin.StorageHandler
	config         *Config
	enabled        bool
	stop           chan struct{}
	wg             *sync.WaitGroup
	ticker         *time.Ticker
}

func (c *RssPlugin) FetchFeed() {
	var storage Storage
	storageBytes, err := c.storageHandler.Load()
	if err != nil {
		log.Printf("could not load storage data: %v", err)
	}
	err = json.Unmarshal(storageBytes, &storage)
	if err != nil {
		log.Printf("could not parse storage data: %v", err)
	}

	fp := gofeed.NewParser()
	for _, url := range c.config.FeedURLs {
		feed, err := fp.ParseURL(url)
		if err != nil {
			log.Printf("error while fetching feed: %v", err)
			continue
		}

		for _, item := range feed.Items {
			log.Printf("Parsing entry %s", item.Title)
			if item.PublishedParsed.After(storage.LastPublished) {
				storage.LastPublished = *item.PublishedParsed
				_ = c.msgHandler.SendMessage(plugin.Message{
					Title:   item.Title,
					Message: item.Link,
				})
			}
		}
	}

	newStorage, err := json.Marshal(storage)
	if err != nil {
		log.Printf("could not marshal storage data: %v", err)
		return
	}
	err = c.storageHandler.Save(newStorage)
	if err != nil {
		log.Printf("could not save storage data: %v", err)
	}
}

func (c *RssPlugin) Enable() error {
	if c.enabled {
		return fmt.Errorf("plugin already enabled")
	}

	c.wg = &sync.WaitGroup{}
	c.stop = make(chan struct{})
	c.ticker = time.NewTicker(time.Duration(c.config.RefreshInterval) * time.Second)
	c.enabled = true

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.stop:
				return
			case <-c.ticker.C:
				c.FetchFeed()
			}
		}
	}()
	return nil
}

func (c *RssPlugin) Disable() error {
	if c.enabled {
		c.enabled = false
		c.ticker.Stop()
		close(c.stop)
		c.wg.Wait()
	} else {
		return fmt.Errorf("plugin already disabled")
	}
	return nil
}

func (c *RssPlugin) GetDisplay(location *url.URL) string {
	var storage Storage

	storageBytes, err := c.storageHandler.Load()
	if err != nil {
		return fmt.Sprintf("Could not load storage data: %v", err)
	}

	err = json.Unmarshal(storageBytes, &storage)
	if err != nil {
		return fmt.Sprintf("Could not parse storage data: %v (%v)", err, storageBytes)
	}

	if storage.LastPublished.IsZero() {
		return fmt.Sprintf("Feed has not been updated as of yet")
	} else {
		return fmt.Sprintf("Last entry was published at %s", storage.LastPublished)
	}
}

func (c *RssPlugin) SetStorageHandler(h plugin.StorageHandler) {
	c.storageHandler = h
}

func (c *RssPlugin) SetMessageHandler(h plugin.MessageHandler) {
	c.msgHandler = h
}

func (c *RssPlugin) DefaultConfig() interface{} {
	return &Config{
		FeedURLs: []string{
			"https://lorem-rss.herokuapp.com/feed",
			"https://xkcd.com/rss.xml",
			"https://news.ycombinator.com/rss",
		},
		RefreshInterval: 3600,
	}
}

func (c *RssPlugin) ValidateAndSetConfig(config interface{}) error {
	c.config = config.(*Config)
	return nil
}

func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &RssPlugin{}
}

func main() {
	panic("this should be built as go plugin")
}
