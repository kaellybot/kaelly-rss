package entities

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type FeedSource struct {
	FeedTypeId string        `gorm:"primaryKey"`
	Url        string        `gorm:"primaryKey"`
	Locale     amqp.Language `gorm:"primaryKey"`
	LastUpdate time.Time     `gorm:"not null; default:current_timestamp"`
}
