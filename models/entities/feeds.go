package entities

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type FeedSource struct {
	FeedTypeID string `gorm:"primaryKey"`
	URL        string
	Game       amqp.Game     `gorm:"primaryKey"`
	Locale     amqp.Language `gorm:"primaryKey"`
	LastUpdate time.Time     `gorm:"not null; default:current_timestamp"`
}
