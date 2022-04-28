package publish

import "time"

type SinglePublishStrategy struct {
	sleepTime time.Duration
}

func (s *SinglePublishStrategy) Wait() {
	time.Sleep(s.sleepTime)
}

func NewSinglePublishStrategy(sleepTime time.Duration) *SinglePublishStrategy {
	return &SinglePublishStrategy{
		sleepTime: sleepTime,
	}
}