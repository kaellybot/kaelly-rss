package models

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	Name               = "KaellyBot"
	RSSParserUserAgent = Name
	RSSParserTimeout   = 60
	RabbitMqAddress    = "amqp://localhost:5672"
	RabbitMqClientId   = "Kaelly-RSS"
)

var (
	RSSUrls = map[amqp.RabbitMQMessage_Language]string{
		amqp.RabbitMQMessage_FR: "https://www.dofus.com/fr/rss/news.xml",
		amqp.RabbitMQMessage_EN: "https://www.dofus.com/en/rss/news.xml",
		amqp.RabbitMQMessage_ES: "https://www.dofus.com/es/rss/news.xml",
	}
)
