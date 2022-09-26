package entity

type Entity struct {
	//index int64
	Data []byte
}

func (e *Entity) DataByte() []byte {
	return e.Data
}

type FactoryEntity struct {
}

func (receiver *FactoryEntity) NewEntity() IEntity {
	return &Entity{}
}
