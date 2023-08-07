package value_object

import (
	"ddd_go_example/internal/app/domain/custom_error"
	"fmt"
	"time"
)

type Date struct {
	year  int
	month int
	day   int
}

func NewDate(year, month, day int) (Date, error) {
	if _, err := toTime(year, month, day); err != nil {
		return Date{}, custom_error.NewBusinessRuleError(custom_error.StatusBadRequest, "invalid date")
	}
	return Date{year: year, month: month, day: day}, nil
}

func toTime(year, month, day int) (time.Time, error) {
	value := fmt.Sprintf("%04d-%02d-%02d", year, month, day)
	return time.Parse(time.DateOnly, value)
}

func (d Date) CalcAge() int {
	thisYear, thisMonth, thisDay := time.Now().Date()
	thisMonthNum := int(thisMonth)
	age := thisYear - d.year
	// 誕生日を迎えていなければ年齢を1つ下げる
	if thisMonthNum < d.month || (thisMonthNum == d.month && thisDay < d.day) {
		age -= 1
	}
	return age
}
