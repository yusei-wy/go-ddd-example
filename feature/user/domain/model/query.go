package model

import "github.com/google/uuid"

type User struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUser(id uuid.UUID, name string) User {
	return User{
		Id:   id,
		Name: name,
	}
}
