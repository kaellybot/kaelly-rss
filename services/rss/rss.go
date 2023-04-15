package rss

import (
	"context"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/models"
	"github.com/kaellybot/kaelly-rss/models/constants"
	"github.com/kaellybot/kaelly-rss/repositories/feedsources"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New(feedSourceRepo feedsources.FeedSourceRepository,
	broker amqp.MessageBrokerInterface) (*RSSServiceImpl, error) {

	fp := gofeed.NewParser()
	fp.UserAgent = constants.RssUserAgent
	return &RSSServiceImpl{
		broker:         broker,
		feedParser:     fp,
		timeout:        time.Duration(viper.GetInt(constants.RssTimeout)) * time.Second,
		feedSourceRepo: feedSourceRepo,
	}, nil
}

func (service *RSSServiceImpl) DispatchNewFeeds() error {
	log.Info().Msgf("Checking feeds...")

	// TODO retrieve feedSources from repo

	var wg sync.WaitGroup
	for _, feedSource := range models.FeedSources {
		wg.Add(1)
		go func(feedSource models.FeedSource) {
			defer wg.Done()
			service.checkFeed(feedSource)
		}(feedSource)
	}

	wg.Wait()
	return nil
}

func (service *RSSServiceImpl) checkFeed(source models.FeedSource) {
	log.Info().
		Str(constants.LogLanguage, source.Language.String()).
		Str(constants.LogUrl, source.Url).
		Msgf("Reading feed source...")
	feed, err := service.readFeed(source.Url)
	if err != nil {
		log.Error().
			Err(err).
			Str(constants.LogLanguage, source.Language.String()).
			Str(constants.LogUrl, source.Url).
			Msgf("Cannot parse URL, source ignored")
		return
	}

	publishedFeeds := 0
	for i := len(feed.Items) - 1; i >= 0; i-- {
		// TODO retrieve new items compared to last time (database access)
		currentFeed := feed.Items[i]
		if currentFeed.PublishedParsed != nil && currentFeed.PublishedParsed.UTC().After(time.Time{}) {
			err := service.publishFeedItem(currentFeed, feed.Copyright, source.Language)
			if err != nil {
				log.Error().Err(err).Msgf("Impossible to publish RSS feed, breaking loop")
				break
			}
			publishedFeeds++
		}
	}

	log.Info().
		Str(constants.LogLanguage, source.Language.String()).
		Str(constants.LogUrl, source.Url).
		Int(constants.LogFeedNumber, publishedFeeds).
		Msgf("Feed(s) read and published")
}

func (service *RSSServiceImpl) readFeed(url string) (*gofeed.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.timeout)
	defer cancel()
	return service.feedParser.ParseURLWithContext(url, ctx)
}

func (service *RSSServiceImpl) publishFeedItem(item *gofeed.Item, source string, language amqp.Language) error {
	msg := models.MapFeedItem(item, source, language)
	return service.broker.Publish(msg, "news", "news.rss", item.GUID)
}
