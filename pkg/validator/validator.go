package validator

import (
	"time"
)

// IsEmpty は文字列が空かどうかを確認します
func IsEmpty(s string) bool {
	return s == ""
}

// IsPastTime は指定された時間が現在より過去かどうかを確認します
func IsPastTime(t time.Time) bool {
	return t.Before(time.Now())
}

// IsFutureTime は指定された時間が現在より未来かどうかを確認します
func IsFutureTime(t time.Time) bool {
	return t.After(time.Now())
}
