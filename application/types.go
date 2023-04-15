package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/services/rss"
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	rss    rss.RSSService
	broker amqp.MessageBrokerInterface
}
