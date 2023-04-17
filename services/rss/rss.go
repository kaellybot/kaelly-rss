package rss

import (
	"context"
	"sort"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/models/constants"
	"github.com/kaellybot/kaelly-rss/models/entities"
	"github.com/kaellybot/kaelly-rss/models/mappers"
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

	feedSources, err := service.feedSourceRepo.GetFeedSources()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, feedSource := range feedSources {
		wg.Add(1)
		go func(feedSource entities.FeedSource) {
			defer wg.Done()
			service.checkFeed(feedSource)
		}(feedSource)
	}

	wg.Wait()
	return nil
}

func (service *RSSServiceImpl) checkFeed(source entities.FeedSource) {
	log.Info().
		Str(constants.LogLanguage, source.Locale.String()).
		Str(constants.LogFeedUrl, source.Url).
		Str(constants.LogFeedType, source.FeedTypeId).
		Msgf("Reading feed source...")

	feed, err := service.readFeed(source.Url)
	if err != nil {
		log.Error().
			Err(err).
			Str(constants.LogLanguage, source.Locale.String()).
			Str(constants.LogFeedType, source.FeedTypeId).
			Str(constants.LogFeedUrl, source.Url).
			Msgf("Cannot parse URL, source ignored")
		return
	}

	publishedFeeds := 0
	lastUpdate := source.LastUpdate
	for _, feedItem := range feed.Items {
		if feedItem.PublishedParsed.UTC().After(lastUpdate.UTC()) {

			err := service.publishFeedItem(feedItem, feed.Copyright, source.Locale)
			if err != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationId, feedItem.GUID).
					Str(constants.LogFeedType, source.FeedTypeId).
					Str(constants.LogFeedUrl, source.Url).
					Str(constants.LogLanguage, source.Locale.String()).
					Str(constants.LogFeedItemId, feedItem.GUID).
					Msgf("Impossible to publish RSS feed, breaking loop")
				break
			}

			source.LastUpdate = feed.PublishedParsed.UTC()
			err = service.feedSourceRepo.Save(source)
			if err != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationId, feedItem.GUID).
					Str(constants.LogFeedType, source.FeedTypeId).
					Str(constants.LogFeedUrl, source.Url).
					Str(constants.LogLanguage, source.Locale.String()).
					Str(constants.LogFeedItemId, feedItem.GUID).
					Msgf("Impossible to update feed source, breaking loop; this feed might be published again next time")
				break
			}

			publishedFeeds++
		}
	}

	log.Info().
		Str(constants.LogLanguage, source.Locale.String()).
		Str(constants.LogFeedType, source.FeedTypeId).
		Str(constants.LogFeedUrl, source.Url).
		Int(constants.LogFeedNumber, publishedFeeds).
		Msgf("Feed(s) read and published")
}

func (service *RSSServiceImpl) readFeed(url string) (*gofeed.Feed, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.timeout)
	defer cancel()
	feed, err := service.feedParser.ParseURLWithContext(url, ctx)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(feed.Items, func(i, j int) bool {
		return feed.Items[i].PublishedParsed.Before(*feed.Items[j].PublishedParsed)
	})

	return feed, nil
}

func (service *RSSServiceImpl) publishFeedItem(item *gofeed.Item, source string, language amqp.Language) error {
	msg := mappers.MapFeedItem(item, source, language)
	return service.broker.Publish(msg, amqp.ExchangeNews, routingkey, item.GUID)
}
