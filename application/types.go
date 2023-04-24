package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/services/feeds"
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	feedService feeds.RSSService
	broker      amqp.MessageBrokerInterface
}
