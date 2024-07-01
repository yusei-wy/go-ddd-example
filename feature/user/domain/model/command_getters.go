// Code generated. DO NOT EDIT.
package model

import (
	"time"
)

func (n UserCommand) Id() UserId {
	return n.id
}

func (n UserCommand) Name() UserName {
	return n.name
}

func (n UserCommand) CreateAt() time.Time {
	return n.createAt
}

func (n UserCommand) UpdateAt() time.Time {
	return n.updateAt
}
