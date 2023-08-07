package user

import (
	"ddd_go_example/internal/app/domain/model/value_object"
)

type UserBirthday = value_object.Date

func NewUserBirthday(year int, month int, day int) (UserBirthday, error) {
	date, err := value_object.NewDate(year, month, day)
	if err != nil {
		return UserBirthday{}, err
	}
	return date, nil
}
