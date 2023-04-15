package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-rss/models/constants"
	"github.com/kaellybot/kaelly-rss/repositories/feedsources"
	"github.com/kaellybot/kaelly-rss/services/rss"
	"github.com/kaellybot/kaelly-rss/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Application, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker, err := amqp.New(constants.RabbitMQClientId, viper.GetString(constants.RabbitMqAddress), nil)
	if err != nil {
		return nil, err
	}

	// repositories
	feedSourcesRepo := feedsources.New(db)

	// services
	rss, err := rss.New(feedSourcesRepo, broker)
	if err != nil {
		return nil, err
	}

	return &Application{rss: rss, broker: broker}, nil
}

func (app *Application) Run() error {
	return app.rss.DispatchNewFeeds()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
