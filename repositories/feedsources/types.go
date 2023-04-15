package feedsources

import (
	"github.com/kaellybot/kaelly-rss/models/entities"
	"github.com/kaellybot/kaelly-rss/utils/databases"
)

type FeedSourceRepository interface {
	GetFeedSources() ([]entities.FeedSource, error)
	Save(feedSource entities.FeedSource) error
}

type FeedSourceRepositoryImpl struct {
	db databases.MySQLConnection
}
