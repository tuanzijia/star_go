package starGo

import (
	"fmt"
	"time"
)

// ToDateTime 转换成日期格式
func ToDateTime(timeVal string) (time.Time, error) {
	if IsEmpty(timeVal) {
		return time.Time{}, fmt.Errorf("输入是空字符串")
	}

	return time.ParseInLocation("2006-01-02 15:04:05", timeVal, time.Local)
}

// ToDate 转换成时间格式
func ToDate(timeVal string) (time.Time, error) {
	if IsEmpty(timeVal) {
		return time.Time{}, fmt.Errorf("输入是空字符串")
	}

	return time.ParseInLocation("2006-01-02", timeVal, time.Local)
}

// ToDateTimeString 转换成时间字符串
func ToDateTimeString(timeVal time.Time) string {
	return timeVal.Local().Format("2006-01-02 15:04:05")
}

// ToDateString 转换成日期字符串
func ToDateString(timeVal time.Time) string {
	return timeVal.Local().Format("2006-01-02")
}
