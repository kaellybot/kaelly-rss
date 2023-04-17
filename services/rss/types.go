package rss

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/repositories/feedsources"
	"github.com/mmcdole/gofeed"
)

const (
	routingkey = "news.rss"
)

type RSSService interface {
	DispatchNewFeeds() error
}

type RSSServiceImpl struct {
	broker         amqp.MessageBrokerInterface
	feedParser     *gofeed.Parser
	timeout        time.Duration
	feedSourceRepo feedsources.FeedSourceRepository
}
