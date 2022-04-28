package subscribe

import (
	"time"
)

type SingleSubscribeStrategy struct {
	sleepTime time.Duration
}

func (s *SingleSubscribeStrategy) Wait() {
	time.Sleep(s.sleepTime)
	return
}

func NewSingleSubscribeStrategy(sleepTime time.Duration) *SingleSubscribeStrategy {
	return &SingleSubscribeStrategy{
		sleepTime: sleepTime,
	}
}
