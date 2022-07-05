package rss

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	mock := NewMock()
	mock.CheckFeeds()

	var called bool
	mock.CheckFeedsFunc = func() {
		called = true
	}

	called = false
	mock.CheckFeeds()
	assert.True(t, called)
}
