package entities

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type FeedSource struct {
	Id         string `gorm:"primaryKey"`
	Language   amqp.Language
	LastUpdate time.Time
}
