package feedsources

import (
	"github.com/kaellybot/kaelly-rss/models/entities"
	"github.com/kaellybot/kaelly-rss/utils/databases"
)

func New(db databases.MySQLConnection) *FeedSourceRepositoryImpl {
	return &FeedSourceRepositoryImpl{db: db}
}

func (repo *FeedSourceRepositoryImpl) GetFeedSources() ([]entities.FeedSource, error) {
	var feedSources []entities.FeedSource
	response := repo.db.GetDB().Model(&entities.FeedSource{}).Find(&feedSources)
	return feedSources, response.Error
}

func (repo *FeedSourceRepositoryImpl) Save(feedSource entities.FeedSource) error {
	return repo.db.GetDB().Save(&feedSource).Error
}
