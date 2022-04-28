package core

type IPublisher interface {
	Start()
	Stop()
	Pub(entity IEntity) error
}
