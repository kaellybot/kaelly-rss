package feeds

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

func New(feedSourceRepo feedsources.Repository, broker amqp.MessageBroker) (*RSSServiceImpl, error) {
	fp := gofeed.NewParser()
	fp.UserAgent = viper.GetString(constants.UserAgent)
	return &RSSServiceImpl{
		broker:         broker,
		feedParser:     fp,
		timeout:        time.Duration(viper.GetInt(constants.RSSTimeout)) * time.Second,
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
		Str(constants.LogFeedURL, source.URL).
		Str(constants.LogFeedType, source.FeedTypeID).
		Msgf("Reading feed source...")

	feed, err := service.readFeed(source.URL)
	if err != nil {
		log.Error().
			Err(err).
			Str(constants.LogLanguage, source.Locale.String()).
			Str(constants.LogFeedType, source.FeedTypeID).
			Str(constants.LogFeedURL, source.URL).
			Msgf("Cannot parse URL, source ignored")
		return
	}

	publishedFeeds := 0
	lastUpdate := source.LastUpdate
	for _, feedItem := range feed.Items {
		if feedItem.PublishedParsed.UTC().After(lastUpdate.UTC()) {
			errPublish := service.publishFeedItem(feedItem, feed.Copyright, source)
			if errPublish != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationID, feedItem.GUID).
					Str(constants.LogFeedType, source.FeedTypeID).
					Str(constants.LogFeedURL, source.URL).
					Str(constants.LogLanguage, source.Locale.String()).
					Str(constants.LogFeedItemID, feedItem.GUID).
					Msgf("Impossible to publish RSS feed, breaking loop")
				break
			}

			source.LastUpdate = feed.PublishedParsed.UTC()
			err = service.feedSourceRepo.Save(source)
			if err != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationID, feedItem.GUID).
					Str(constants.LogFeedType, source.FeedTypeID).
					Str(constants.LogFeedURL, source.URL).
					Str(constants.LogLanguage, source.Locale.String()).
					Str(constants.LogFeedItemID, feedItem.GUID).
					Msgf("Impossible to update feed source, breaking loop; this feed might be published again next time")
				break
			}

			publishedFeeds++
		}
	}

	log.Info().
		Str(constants.LogLanguage, source.Locale.String()).
		Str(constants.LogFeedType, source.FeedTypeID).
		Str(constants.LogFeedURL, source.URL).
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

func (service *RSSServiceImpl) publishFeedItem(item *gofeed.Item, source string,
	feedSource entities.FeedSource) error {
	msg := mappers.MapFeedItem(item, source, feedSource)
	return service.broker.Emit(msg, amqp.ExchangeNews, routingkey, item.GUID)
}
