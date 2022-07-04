package rss

import (
	"context"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/models"
	"github.com/mmcdole/gofeed"
	"github.com/rs/zerolog/log"
)

type RSSServiceInterface interface {
	CheckFeeds()
}

type RSSService struct {
	broker     amqp.MessageBrokerInterface
	feedParser *gofeed.Parser
	timeout    time.Duration
}

func New(broker amqp.MessageBrokerInterface, timeout int) (*RSSService, error) {
	fp := gofeed.NewParser()
	fp.UserAgent = models.RSSParserUserAgent
	return &RSSService{
		broker:     broker,
		feedParser: fp,
		timeout:    time.Duration(timeout) * time.Second,
	}, nil
}

func (service *RSSService) CheckFeeds() {
	log.Info().Msgf("Checking feeds...")

	var wg sync.WaitGroup
	for language, url := range models.RSSUrls {
		wg.Add(1)
		go func(language amqp.RabbitMQMessage_Language, url string) {
			defer wg.Done()
			service.checkFeeds(language, url)
		}(language, url)
	}

	wg.Wait()
}

func (service *RSSService) checkFeeds(language amqp.RabbitMQMessage_Language, url string) {
	log.Info().
		Interface(models.LogLanguage, language).
		Str(models.LogUrl, url).
		Msgf("Reading feed source...")

	ctx, cancel := context.WithTimeout(context.Background(), service.timeout)
	defer cancel()

	feed, err := service.feedParser.ParseURLWithContext(url, ctx)
	if err != nil {
		log.Error().Err(err).
			Interface(models.LogLanguage, language).
			Str(models.LogUrl, url).
			Msgf("Cannot parse URL, source ignored")
		return
	}

	publishedFeeds := 0
	for i := len(feed.Items) - 1; i >= 0; i-- {
		// TODO retrieve new items compared to last time (database access)
		currentFeed := feed.Items[i]
		if currentFeed.PublishedParsed != nil && currentFeed.PublishedParsed.UTC().After(time.Time{}) {
			msg := models.MapFeed(currentFeed, language)
			err := service.broker.Publish(msg, "news", "news.rss", currentFeed.GUID)
			if err != nil {
				log.Error().Err(err).Msgf("Impossible to publish RSS feed, breaking loop")
				break
			}
			publishedFeeds++
		}
	}

	log.Info().
		Interface(models.LogLanguage, language).
		Int(models.LogFeedNumber, publishedFeeds).
		Msgf("Feed(s) read and published")
}
