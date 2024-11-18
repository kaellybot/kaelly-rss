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
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress))
	db := databases.New()

	// repositories
	feedSourcesRepo := feedsources.New(db)

	// services
	feedService, err := feeds.New(feedSourcesRepo, broker)
	if err != nil {
		return nil, err
	}

	return &Impl{
		feedService: feedService,
		broker:      broker,
		db:          db,
	}, nil
}

func (app *Impl) Run() error {
	if err := app.db.Run(); err != nil {
		return err
	}

	if err := app.broker.Run(); err != nil {
		return err
	}

	return app.feedService.DispatchNewFeeds()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	app.db.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
