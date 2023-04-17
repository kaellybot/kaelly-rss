package mappers

import (
	"regexp"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/mmcdole/gofeed"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	imageUrlRegex, _ = regexp.Compile("<img.+src=\"(.*\\.jpg)\".+>")
)

func MapFeedItem(item *gofeed.Item, source string, language amqp.Language) *amqp.RabbitMQMessage {
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
