package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/models/constants"
	"github.com/kaellybot/kaelly-rss/repositories/feedsources"
	"github.com/kaellybot/kaelly-rss/services/feeds"
	"github.com/kaellybot/kaelly-rss/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress))

	// repositories
	feedSourcesRepo := feedsources.New(db)

	// services
	feedService, err := feeds.New(feedSourcesRepo, broker)
	if err != nil {
		return nil, err
	}

	return &Impl{feedService: feedService, broker: broker}, nil
}

func (app *Impl) Run() error {
	if err := app.broker.Run(); err != nil {
		return err
	}

	return app.feedService.DispatchNewFeeds()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
