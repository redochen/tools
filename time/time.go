package time

import (
	"fmt"
	"strings"
	"time"

	CcStr "github.com/redochen/tools/string"
)

const (
	DefaultTimeZone       = ""
	DefaultDateFormat     = "yyyy-MM-dd"
	DefaultTimeFormat     = "HH:mm:ss"
	DefaultDateTimeFormat = DefaultDateFormat + " " + DefaultTimeFormat
)

// DaysOfMonth 每月的天数
var DaysOfMonth = map[time.Month]int{
	time.January:   31,
	time.February:  28,
	time.March:     31,
	time.April:     30,
	time.May:       31,
	time.June:      30,
	time.July:      31,
	time.August:    31,
	time.September: 30,
	time.October:   31,
	time.November:  30,
	time.December:  31,
}

// DateDifference 日期差
type DateDifference struct {
	Years  int
	Months int
	Days   int
}

// IsLeapYear 是否为润年
func IsLeapYear(year int) bool {
	ret := false
	if year%4 == 0 {
		if year%100 != 0 {
			ret = true
		} else if year%400 == 0 {
			ret = true
		}
	}
	return ret
}

// CalculateDateDifference 计算日期差
func CalculateDateDifference(start, end time.Time) *DateDifference {
	var borrowed = false
	var daysBorrowed = 0
	diff := &DateDifference{}

	// Days difference
	if end.Day() >= start.Day() {
		diff.Days = end.Day() - start.Day()
	} else {
		daysBorrowed = DaysOfMonth[end.Month()-1]
		if IsLeapYear(end.Year()) &&
			end.Month() == time.March {
			daysBorrowed++ // February in leap year is 29 days
		}
		diff.Days = end.Day() + daysBorrowed - start.Day()
		borrowed = true
	}

	// Month Difference
	endMonth := end.Month()
	if borrowed == true {
		if endMonth == time.January {
			endMonth = time.December
		} else {
			endMonth--
		}
		borrowed = false
	}

	if endMonth >= start.Month() {
		diff.Months = int(endMonth) - int(start.Month())
	} else {
		diff.Months = int(endMonth) + 12 - int(start.Month())
		borrowed = true
	}

	// Year difference
	if borrowed {
		diff.Years = end.Year() - 1 - start.Year()
	} else {
		diff.Years = end.Year() - start.Year()
	}

	return diff
}

// RemoveDateSeparator 去掉日期中的连接符-/和空格
func RemoveDateSeparator(date string) string {
	if len(date) > 0 {
		separators := []string{"-", "/", " "}
		for _, sep := range separators {
			date = strings.Replace(date, sep, "", -1)
		}
	}
	return date
}

// GetDateSeparator 获取日期分隔符
func GetDateSeparator(date string) string {
	if len(date) <= 0 {
		return ""
	}

	separators := []string{"-", "/", " ", "."}
	for _, sep := range separators {
		if strings.Contains(date, sep) {
			return sep
		}
	}

	return ""
}

// AddDateSeparator 将YYYYMMDD格式的日期转换为YYYY-MM-DD格式或者YYYY-MM格式
func AddDateSeparator(date, separator string, incDay bool) string {
	temp := RemoveDateSeparator(date)
	if len(temp) < 8 {
		return date
	}

	if incDay {
		return fmt.Sprintf("%s%s%s%s%s", temp[:4], separator, temp[4:6], separator, temp[6:8])
	}

	return fmt.Sprintf("%s%s%s", temp[:4], separator, temp[4:6])
}

// RemoveTimeFromDateTime 将日期时间字符串中的时间信息去掉
func RemoveTimeFromDateTime(date string) string {
	separator := GetDateSeparator(date)
	return AddDateSeparator(date, separator, true)
}

// GetShortDate 获取短日期字符串
func GetShortDate(date string) string {
	temp := RemoveDateSeparator(date)
	if len(temp) < 8 {
		return date
	}
	return fmt.Sprintf("%s%s", temp[4:6], temp[6:8])
}

// NowTime 获取当前时间（默认时区）
func NowTime() time.Time {
	return NowTimeEx(DefaultTimeZone)
}

// NowTimeEx 获取当前时间（指定时区）
func NowTimeEx(timeZone string) time.Time {
	nowTime := NowTimeTextEx(timeZone, DefaultDateTimeFormat)
	return CcStr.ParseTime(nowTime, DefaultDateTimeFormat, false)
}

// NowTimeText 获取当前时间（默认时区）
func NowTimeText(format string) string {
	return NowTimeTextEx(DefaultTimeZone, format)
}

// NowTimeTextEx 获取当前时间（指定时区）
func NowTimeTextEx(timeZone, format string) string {
	now := time.Now()
	if timeZone != "" {
		location, err := time.LoadLocation(timeZone)
		if err == nil {
			now = now.In(location)
		}
	}

	//默认格式：yyyy-MM-dd HH:mm:ss
	if format == "" {
		format = DefaultDateTimeFormat
	}
	return CcStr.FormatTime(now, format)
}

// ToInt64 将日期时间格式化字符串后再转换成64位整数
func ToInt64(t time.Time, format string) int64 {
	s := CcStr.FormatTime(t, format)
	return CcStr.ParseInt64(s)
}

// ToDayStart 将日期时间转换成 yyyy-MM-dd 00:00:00:000
func ToDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// ToDayEnd 将日期时间转换成 yyyy-MM-dd 23:59:59:999
func ToDayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location())
}
