package model

import (
	customerror "go_ddd_example/share/custom_error"
	"go_ddd_example/share/domain/model/valueobject"

	"github.com/google/uuid"
)

type UserID struct {
	valueobject.ValueObject[uuid.UUID]
}

func NewUserID() UserID {
	return UserID{valueobject.NewValueObject(uuid.New())}
}

func ParseUserID(userID string) (UserID, customerror.ModelError) {
	u, err := uuid.Parse(userID)
	if err != nil {
		return UserID{valueobject.NewValueObject(uuid.Nil)}, customerror.NewModelErrorWithMessage(err, "Invalid user id")
	}

	return UserID{valueobject.NewValueObject(u)}, nil
}

type UserName struct {
	valueobject.ValueObject[string]
}

func ParseUserName(name string) (UserName, customerror.ModelError) {
	if name == "" {
		return UserName{valueobject.NewValueObject("")}, customerror.NewModelErrorWithMessage(nil, "Name is required")
	}

	return UserName{valueobject.NewValueObject(name)}, nil
}
