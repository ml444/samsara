package entity

import "github.com/ml444/samsara/core"

type Entity struct {
	//index int64
	Data  []byte
}

func (e *Entity) DataByte() []byte {
	return e.Data
}

type FactoryEntity struct {
}

func (receiver *FactoryEntity) NewEntity() core.IEntity {
	return &Entity{}
}
