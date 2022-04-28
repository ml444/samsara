package publish

import (
	"errors"
	"github.com/ml444/samsara/core"
)

type MultiPublisher struct {
	barrier  core.IPublishBarrier
	strategy core.IPublisherStrategy
	isDone   bool
}

func NewMultiPublisher(scheduler core.IScheduler, strategy core.IPublisherStrategy) *MultiPublisher {
	return &MultiPublisher{
		barrier:  core.NewMultiPublishBarrier(scheduler, strategy),
		strategy: strategy,
	}
}

func (p *MultiPublisher) Start() {

}
func (p *MultiPublisher) Stop() {
	p.isDone = true
}
func (p *MultiPublisher) Pub(entity core.IEntity) error {
	if p.isDone {
		return errors.New("shutdown")
	}
	seq := p.barrier.Next()
	//println(seq)
	p.barrier.Commit(seq, entity)
	return nil
}
func (p *MultiPublisher) PubNoWait(entity core.IEntity) error {
	if p.isDone {
		return errors.New("shutdown")
	}
	seq, err := p.barrier.TryNext()
	if err != nil {
		return err
	}
	p.barrier.Commit(seq, entity)
	return nil
}
