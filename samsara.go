package samsara

import (
	"github.com/ml444/samsara/core"
	"github.com/ml444/samsara/utils"
	"time"
)

type Samsara struct {
	ringBuffer *core.RingBuffer
	scheduler  core.IScheduler
	//subscriberRepository *subscribe.Repository

	subscriberList []core.ISubscriber
	publisher      core.IPublisher
	isDone         *utils.AtomicBool
	//exceptionHandler   *ExceptionHandler
	once utils.Once
}

func NewSamsara(ringBufferSize int64, eventFactory core.IEntityFactory) *Samsara {
	ringBuffer := core.NewRingBuffer(ringBufferSize, eventFactory)
	scheduler := core.NewScheduler(ringBuffer)
	return &Samsara{
		ringBuffer: ringBuffer,
		scheduler:  scheduler,
		//subscriberRepository: &subscribe.Repository{},
		isDone: &utils.AtomicBool{},
	}
}

func (s *Samsara) GetScheduler() core.IScheduler {
	return s.scheduler
}

func (s *Samsara) SetPublisher(publisher core.IPublisher) {
	s.publisher = publisher
}

func (s *Samsara) SetSubscriber(subscriber core.ISubscriber) {
	s.scheduler.AddSequences(subscriber.GetSequence())
	s.subscriberList = append(s.subscriberList, subscriber)
}

func (s *Samsara) Start() {
	subscribers := s.subscriberList
	//sequences := s.scheduler.GetConsumerSequences()
	//s.scheduler.InitConsumerSequences(sequences...)

	s.once.Do(func() {
		for _, subscriberInst := range subscribers {
			go subscriberInst.Start()
		}

		go s.publisher.Start()
	})
	s.isDone.Set(false)
	return
}

func (s *Samsara) StopPublisher() {
	s.publisher.Stop()
}

func (s *Samsara) StopSubscribers() {
	for _, subscriberInst := range s.subscriberList {
		subscriberInst.Stop()
	}
}

func (s *Samsara) Shutdown() {
	s.StopPublisher()
	for s.HasBlocking() {
		time.Sleep(time.Millisecond * 1000)
	}
	s.StopSubscribers()
	s.isDone.Set(true)
}

func (s *Samsara) HasBlocking() bool {
	cursor := s.scheduler.GetCursor()
	for _, subscriber := range s.subscriberList {
		println(subscriber, cursor.Get(),subscriber.GetSequence().Get())
		if cursor.Get() != subscriber.GetSequence().Get() {
			return true
		}
	}
	return false
}

func (s *Samsara) IsDone() bool {
	return s.isDone.Get()
}

func (s *Samsara) GetRingBuffer() *core.RingBuffer {
	return s.ringBuffer
}

func (s *Samsara) GetCursor() int64 {
	return s.scheduler.GetCursor().Get()
}

func (s *Samsara) GetBufferSize() int64 {
	return s.scheduler.GetBufferSize()
}
func (s *Samsara) Publish(entity core.IEntity) error {
	return s.publisher.Pub(entity)
}

func (s *Samsara) Get(sequence int64) core.IEntity {
	return s.ringBuffer.GetEntity(sequence)

}
