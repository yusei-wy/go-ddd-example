package value_object

import (
	"testing"
	"time"
)

func TestDate_CalcAge(t *testing.T) {
	now := time.Now()
	year, month, day := now.Date()
	monthNum := int(month)

	tests := []struct {
		name   string
		fields Date
		want   int
	}{
		{
			name: "0 years old",
			fields: Date{
				year:  year,
				month: monthNum,
				day:   day,
			},
			want: 0,
		},
		{
			name: "18 years old no.1",
			fields: Date{
				year:  year - 18,
				month: monthNum,
				day:   day,
			},
			want: 18,
		},
		{
			name: "18 years old no.2",
			fields: Date{
				year:  year - 18,
				month: monthNum - 11,
				day:   day - 1,
			},
			want: 18,
		},
		{
			name: "17 years old no.1",
			fields: Date{
				year:  year - 18,
				month: monthNum + 1,
				day:   day,
			},
			want: 17,
		},
		{
			name: "17 years old no.2",
			fields: Date{
				year:  year - 18,
				month: monthNum,
				day:   day + 1,
			},
			want: 17,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := Date{
				year:  tt.fields.year,
				month: tt.fields.month,
				day:   tt.fields.day,
			}
			if got := d.CalcAge(); got != tt.want {
				t.Errorf("Date.CalcAge() = %v, want %v", got, tt.want)
			}
		})
	}
}
