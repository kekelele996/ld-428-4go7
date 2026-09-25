package util

import "time"

// FormatDate 日期。
func FormatDate(t time.Time) string { return t.Format("2006-01-02") }

// FormatDateTime 时间。
func FormatDateTime(t time.Time) string { return t.Format("2006-01-02T15:04:05Z07:00") }
