package models

import (
	"regexp"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/mmcdole/gofeed"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedSource struct {
	Type     FeedSourceType
	Language amqp.RabbitMQMessage_Language
	Url      string
}

type FeedSourceType string

const (
	Changelog FeedSourceType = "changelog"
	Devblog   FeedSourceType = "devblog"
	News      FeedSourceType = "news"
)

var (
	FeedSources = []FeedSource{
		{
			Type:     Changelog,
			Language: amqp.RabbitMQMessage_FR,
			Url:      "https://www.dofus.com/fr/rss/changelog.xml",
		},
		{
			Type:     Devblog,
			Language: amqp.RabbitMQMessage_FR,
			Url:      "https://www.dofus.com/fr/rss/devblog.xml",
		},
		{
			Type:     News,
			Language: amqp.RabbitMQMessage_FR,
			Url:      "https://www.dofus.com/fr/rss/news.xml",
		},
		{
			Type:     Changelog,
			Language: amqp.RabbitMQMessage_EN,
			Url:      "https://www.dofus.com/en/rss/changelog.xml",
		},
		{
			Type:     Devblog,
			Language: amqp.RabbitMQMessage_EN,
			Url:      "https://www.dofus.com/en/rss/devblog.xml",
		},
		{
			Type:     News,
			Language: amqp.RabbitMQMessage_EN,
			Url:      "https://www.dofus.com/en/rss/news.xml",
		},
		{
			Type:     Changelog,
			Language: amqp.RabbitMQMessage_ES,
			Url:      "https://www.dofus.com/es/rss/changelog.xml",
		},
		{
			Type:     Devblog,
			Language: amqp.RabbitMQMessage_ES,
			Url:      "https://www.dofus.com/es/rss/devblog.xml",
		},
		{
			Type:     News,
			Language: amqp.RabbitMQMessage_ES,
			Url:      "https://www.dofus.com/es/rss/news.xml",
		},
	}

	imageUrlRegex, _ = regexp.Compile("<img.+src=\"(.*\\.jpg)\".+>")
)

func MapFeedItem(item *gofeed.Item, source string, language amqp.RabbitMQMessage_Language) *amqp.RabbitMQMessage {
	var iconUrl string
	if item.Image != nil {
		iconUrl = item.Image.URL
	} else if matches := imageUrlRegex.FindStringSubmatch(item.Description); matches != nil && len(matches) >= 2 {
		iconUrl = matches[1]
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_RSS,
		Language: language,
		NewsRSSMessage: &amqp.NewsRSSMessage{
			Title:      item.Title,
			AuthorName: source,
			Url:        item.Link,
			IconUrl:    iconUrl,
			Date:       timestamppb.New(*item.PublishedParsed),
		},
	}
}
