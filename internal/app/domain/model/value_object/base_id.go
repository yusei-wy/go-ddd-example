package value_object

import "github.com/google/uuid"

type BaseId struct {
	value uuid.UUID
}

func NewBaseId() BaseId {
	return BaseId{value: uuid.New()}
}
