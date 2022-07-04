package models

import (
	"regexp"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/mmcdole/gofeed"
)

var (
	ImageUrlRegex, _ = regexp.Compile("<img.+src=\"(.*\\.jpg)\".+>")
)

func MapFeed(feed *gofeed.Item, language amqp.RabbitMQMessage_Language) *amqp.RabbitMQMessage {
	var iconUrl string
	if feed.Image != nil {
		iconUrl = feed.Image.URL
	} else if matches := ImageUrlRegex.FindStringSubmatch(feed.Description); matches != nil && len(matches) >= 2 {
		iconUrl = matches[1]
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_RSS,
		Language: language,
		NewsRSSMessage: &amqp.NewsRSSMessage{
			Title:   feed.Title,
			Url:     feed.Link,
			IconUrl: iconUrl,
			Date:    feed.PublishedParsed.UnixMilli(),
		},
	}
}
