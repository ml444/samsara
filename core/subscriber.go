package core

type ISubscriber interface {
	GetSequence() *Sequence
	Start()
	Stop()
	//MarkAsUsedInBarrier()
	//IsRunning()

	Handler(entity IEntity)
}
