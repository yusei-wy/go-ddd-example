package model

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewUser(id uuid.UUID, name string) User {
	return User{
		ID:   id,
		Name: name,
	}
}
