package main

import (
	"fmt"
	"net/http"

	"github.com/kaellybot/kaelly-rss/application"
	"github.com/kaellybot/kaelly-rss/models/constants"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	initConfig()
	initLog()
	initMetrics()
}

func initConfig() {
	viper.SetConfigFile(constants.ConfigFileName)

	for configName, defaultValue := range constants.DefaultConfigValues {
		viper.SetDefault(configName, defaultValue)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal().Err(err).Str(constants.LogFileName, constants.ConfigFileName).Msgf("Failed to read config, shutting down.")
	}
}

func initLog() {
	zerolog.SetGlobalLevel(constants.LogLevelFallback)
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	log.Logger = log.With().Caller().Logger()

	logLevel, err := zerolog.ParseLevel(viper.GetString(constants.LogLevel))
	if err != nil {
		log.Warn().Err(err).Msgf("Log level not set, continue with %s...", constants.LogLevelFallback)
	} else {
		zerolog.SetGlobalLevel(logLevel)
		log.Debug().Msgf("Logger level set to '%s'", logLevel)
	}
}

func initMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(fmt.Sprintf(":%v", viper.GetInt(constants.MetricPort)), nil)
}

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to instantiate application")
	}

	app.Run()
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to run application")
	}

	log.Info().Msgf("Gracefully shutting down %s...", constants.InternalName)
	app.Shutdown()
}
