package rss

import "github.com/rs/zerolog/log"

type RSSServiceMock struct {
	CheckFeedsFunc func()
}

func (mock *RSSServiceMock) CheckFeeds() {
	if mock.CheckFeedsFunc != nil {
		mock.CheckFeeds()
		return
	}

	log.Warn().Msgf("No mock provided")
}
