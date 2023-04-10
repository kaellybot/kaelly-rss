package application

import (
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/services/rss"
)

var (
	ErrCannotInstanciateApp = errors.New("Cannot instanciate application")
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	rss    rss.RSSServiceInterface
	broker amqp.MessageBrokerInterface
}
