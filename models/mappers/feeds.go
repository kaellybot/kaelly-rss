package mappers

import (
	"regexp"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/mmcdole/gofeed"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	imageURLRegexExpectedGroup = 2
)

var (
	imageURLRegex = regexp.MustCompile("<img.+src=\"(.*\\.jpg)\".+>")
)

func MapFeedItem(item *gofeed.Item, source, feedType string, language amqp.Language) *amqp.RabbitMQMessage {
	var iconURL string
	if item.Image != nil {
		iconURL = item.Image.URL
	} else if matches := imageURLRegex.FindStringSubmatch(item.Description); len(matches) >= imageURLRegexExpectedGroup {
		iconURL = matches[1]
	}

	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_RSS,
		Language: language,
		NewsRSSMessage: &amqp.NewsRSSMessage{
			Title:      item.Title,
			AuthorName: source,
			Url:        item.Link,
			IconUrl:    iconURL,
			Date:       timestamppb.New(*item.PublishedParsed),
			Type:       feedType,
		},
	}
}
