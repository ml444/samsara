package publish

import (
	"github.com/ml444/samsara/entity"
)

type IPublisher interface {
	Init()
	Pause()
	Pub(entity entity.IEntity) error
}
