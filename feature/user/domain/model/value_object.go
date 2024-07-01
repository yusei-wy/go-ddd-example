package model

import (
	customerror "go_ddd_example/share/custom_error"
	"go_ddd_example/share/domain/model/value_object"

	"github.com/google/uuid"
)

type UserId struct {
	value_object.ValueObject[uuid.UUID]
}

func NewUserId() UserId {
	return UserId{value_object.NewValueObject(uuid.New())}
}

func ParseUserId(userId string) (UserId, customerror.ModelError) {
	u, err := uuid.Parse(userId)
	if err != nil {
		return UserId{value_object.NewValueObject(uuid.Nil)}, customerror.NewModelErrorWithMessage(err, "Invalid user id")
	}

	return UserId{value_object.NewValueObject(u)}, nil
}

type UserName struct {
	value_object.ValueObject[string]
}

func ParseUserName(name string) (UserName, customerror.ModelError) {
	if name == "" {
		return UserName{value_object.NewValueObject("")}, customerror.NewModelErrorWithMessage(nil, "Name is required")
	}

	return UserName{value_object.NewValueObject(name)}, nil
}
