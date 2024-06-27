package entity

import "github.com/google/uuid"

type ID = uuid.UUID

func NewID() ID {
	return ID(uuid.New())
}

func ParseID(stringId string) (ID, error) {
	id, err := uuid.Parse(stringId)
	return ID(id), err
}
