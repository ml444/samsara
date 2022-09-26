package publish

import (
	"errors"
	"github.com/ml444/samsara/entity"
	"github.com/ml444/samsara/internal"
)

type Producer struct {
	barrier internal.IPublishBarrier
	//strategy internal.IPublisherStrategy
	isPause bool
}

func NewProducer(barrier internal.IPublishBarrier) *Producer {
	return &Producer{
		barrier: barrier,
		//strategy: strategy,
	}
}

func (p *Producer) Init() {

}
func (p *Producer) Pause() {
	p.isPause = true
}
func (p *Producer) Pub(entity entity.IEntity) error {
	if p.isPause {
		return errors.New("this producer has paused")
	}
	seq := p.barrier.Next()
	p.barrier.Commit(seq, entity)
	return nil
}
func (p *Producer) PubNoWait(entity entity.IEntity) error {
	if p.isPause {
		return errors.New("this producer has paused")
	}
	seq, err := p.barrier.TryNext()
	if err != nil {
		return err
	}
	p.barrier.Commit(seq, entity)
	return nil
}
