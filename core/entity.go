package core

type IEntityFactory interface {
	NewEntity() IEntity
}

type IEntity interface {
	DataByte() []byte
}
