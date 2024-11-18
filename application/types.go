package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/services/feeds"
	"github.com/kaellybot/kaelly-rss/utils/databases"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	feedService feeds.RSSService
	broker      amqp.MessageBroker
	db          databases.MySQLConnection
}
