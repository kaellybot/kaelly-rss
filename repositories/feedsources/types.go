package feedsources

import (
	"github.com/kaellybot/kaelly-rss/models/entities"
	"github.com/kaellybot/kaelly-rss/utils/databases"
)

type Repository interface {
	GetFeedSources() ([]entities.FeedSource, error)
	Save(feedSource entities.FeedSource) error
}

type Impl struct {
	db databases.MySQLConnection
}
