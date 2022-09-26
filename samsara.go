package samsara

import (
	"github.com/ml444/samsara/entity"
	"github.com/ml444/samsara/internal"
	"github.com/ml444/samsara/publish"
	"github.com/ml444/samsara/subscribe"
	"github.com/ml444/samsara/utils"
	"time"
)

type Samsara struct {
	ringBuffer     *internal.RingBuffer
	scheduler      internal.IScheduler
	publisherList  []publish.IPublisher
	subscriberList []subscribe.ISubscriber
	isDone         *utils.AtomicBool
	once           utils.Once
}

func NewSamsara(ringBufferSize int64, eventFactory entity.IEntityFactory) *Samsara {
	ringBuffer := internal.NewRingBuffer(ringBufferSize, eventFactory)
	scheduler := internal.NewScheduler(ringBuffer)
	return &Samsara{
		ringBuffer: ringBuffer,
		scheduler:  scheduler,
		isDone:     &utils.AtomicBool{},
	}
}

func (s *Samsara) GetScheduler() internal.IScheduler {
	return s.scheduler
}

func (s *Samsara) AddPublisher(publisher publish.IPublisher) {
	s.publisherList = append(s.publisherList, publisher)
}

func (s *Samsara) AddSubscriber(subscriber subscribe.ISubscriber) {
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
	})
	s.isDone.Set(false)
	return
}

func (s *Samsara) PausePublishers() {
	for _, publisher := range s.publisherList {
		publisher.Pause()
	}
}
func (s *Samsara) StopSubscribers() {
	for _, subscriberInst := range s.subscriberList {
		subscriberInst.Stop()
	}
}

func (s *Samsara) Shutdown() {
	s.PausePublishers()
	for s.HasBlocking() {
		time.Sleep(time.Millisecond * 100)
	}
	s.StopSubscribers()
	s.isDone.Set(true)
}

func (s *Samsara) HasBlocking() bool {
	cursor := s.scheduler.GetCursor()
	for _, subscriber := range s.subscriberList {
		if cursor.Get() != subscriber.GetSequence().Get() {
			println("blocking:", cursor.Get(), subscriber.GetSequence().Get())
			return true
		}
	}
	return false
}

func (s *Samsara) IsDone() bool {
	return s.isDone.Get()
}

func (s *Samsara) GetRingBuffer() *internal.RingBuffer {
	return s.ringBuffer
}

func (s *Samsara) GetCursor() int64 {
	return s.scheduler.GetCursor().Get()
}

func (s *Samsara) GetBufferSize() int64 {
	return s.scheduler.GetBufferSize()
}

func (s *Samsara) Get(sequence int64) entity.IEntity {
	return s.ringBuffer.GetEntity(sequence)
}

func (s *Samsara) NewSinglePublisher(strategy internal.IPublisherStrategy) publish.IPublisher {
	producer := publish.NewProducer(internal.NewSinglePublishBarrier(s.scheduler, strategy))
	s.AddPublisher(producer)
	return producer
}
func (s *Samsara) NewMultiPublisher(strategy internal.IPublisherStrategy) publish.IPublisher {
	producer := publish.NewProducer(internal.NewMultiPublishBarrier(s.scheduler, strategy))
	s.AddPublisher(producer)
	return producer
}

func (s *Samsara) NewSimpleSubscriber(strategy internal.ISubscriberStrategy, handler func(entity entity.IEntity)) subscribe.ISubscriber {
	consumer := subscribe.NewSimpleSubscriber(s.scheduler, strategy, handler)
	s.AddSubscriber(consumer)
	return consumer
}

func NewPublishStrategy(duration time.Duration) internal.IPublisherStrategy {
	return publish.NewSinglePublishStrategy(duration)
}

func NewSubscribeStrategy(duration time.Duration) internal.ISubscriberStrategy {
	return subscribe.NewSingleSubscribeStrategy(duration)
}
