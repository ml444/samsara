package core

type IPublisherStrategy interface {
	Wait()
}
type ISubscriberStrategy interface {
	Wait()
	//WaitFor(sequence int64, cursor *Sequence, barrier ISubscribeBarrier) int64
}
