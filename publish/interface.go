package publish

type IPublisher interface {
	Init()
	Pause()
	Pub(entity interface{}) error
}
