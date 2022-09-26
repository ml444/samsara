package subscribe

import "github.com/ml444/samsara/internal"

type ISubscriber interface {
	GetSequence() *internal.Sequence
	Start()
	Stop()
}
