package value_user

import (
	"ddd_go_example/internal/app/domain/custom_errors"
	"net/http"
)

type UserName struct {
	value string
}

func NewUserName(value string) (UserName, error) {
	if value == "" {
		return UserName{}, custom_errors.NewBusinessRuleError(http.StatusBadRequest, "user name is empty")
	}
	return UserName{value: value}, nil
}
