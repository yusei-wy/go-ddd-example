package user

import (
	"ddd_go_example/internal/app/domain/custom_error"
)

type UserName struct {
	value string
}

func NewUserName(value string) (UserName, error) {
	if value == "" {
		return UserName{}, custom_error.NewBusinessRuleError(custom_error.StatusBadRequest, "user name is empty")
	}
	return UserName{value: value}, nil
}
