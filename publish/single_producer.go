package publish

import (
	"errors"
	"github.com/ml444/samsara/core"
)

type SinglePublisher struct {
	barrier  core.IPublishBarrier
	strategy core.IPublisherStrategy
	isDone   bool
}

func NewSinglePublisher(scheduler core.IScheduler, strategy core.IPublisherStrategy) *SinglePublisher {
	return &SinglePublisher{
		barrier:  core.NewSinglePublishBarrier(scheduler, strategy),
		strategy: strategy,
	}
}

func (p *SinglePublisher) Start() {

}
func (p *SinglePublisher) Stop() {
	p.isDone = true
}
func (p *SinglePublisher) Pub(entity core.IEntity) error {
	if p.isDone {
		return errors.New("shutdown")
	}
	seq := p.barrier.Next()
	p.barrier.Commit(seq, entity)
	return nil
}
