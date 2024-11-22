package feedsources

import (
	"github.com/kaellybot/kaelly-rss/models/entities"
	"github.com/kaellybot/kaelly-rss/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetFeedSources() ([]entities.FeedSource, error) {
	var feedSources []entities.FeedSource
	response := repo.db.GetDB().Model(&entities.FeedSource{}).Find(&feedSources)
	return feedSources, response.Error
}

func (repo *Impl) Save(feedSource entities.FeedSource) error {
	return repo.db.GetDB().
		Model(&feedSource).
		Where("feed_type_id = ? AND game = ? AND locale = ?",
			feedSource.FeedTypeID, feedSource.Game, feedSource.Locale).
		Update("last_update", feedSource.LastUpdate).
		Error
}
