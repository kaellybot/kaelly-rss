package rss

import "github.com/rs/zerolog/log"

type RSSServiceMock struct {
	CheckFeedsFunc func()
}

func NewMock() *RSSServiceMock {
	return &RSSServiceMock{}
}

func (mock *RSSServiceMock) CheckFeeds() {
	if mock.CheckFeedsFunc != nil {
		mock.CheckFeedsFunc()
		return
	}

	log.Warn().Msgf("No mock provided")
}
