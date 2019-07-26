package time

import (
	"fmt"
	. "github.com/redochen/tools/string"
	"strings"
	"time"
)

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

var (
	CcTime = NewTimeHelper()
)

//time帮助类
type TimeHelper struct {
}

//获取一个新的StringHelper实例
func NewTimeHelper() *TimeHelper {
	return &TimeHelper{}
}

//日期差
type DateDifference struct {
	Years  int
	Months int
	Days   int
}

//是否为润年
func (h *TimeHelper) IsLeapYear(year int) bool {
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

//计算日期差
func (h *TimeHelper) CalculateDateDifference(start, end time.Time) *DateDifference {
	var borrowed = false
	var daysBorrowed = 0
	diff := &DateDifference{}

	// Days difference
	if end.Day() >= start.Day() {
		diff.Days = end.Day() - start.Day()
	} else {
		daysBorrowed = DaysOfMonth[end.Month()-1]
		if h.IsLeapYear(end.Year()) &&
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

//去掉日期中的连接符-/和空格
func (h *TimeHelper) RemoveDateSeparator(date string) string {
	if len(date) > 0 {
		separators := []string{"-", "/", " "}
		for _, sep := range separators {
			date = strings.Replace(date, sep, "", -1)
		}
	}
	return date
}

//获取日期分隔符
func (h *TimeHelper) GetDateSeparator(date string) string {
	if len(date) > 0 {
		separators := []string{"-", "/", " ", "."}
		for _, sep := range separators {
			if strings.Contains(date, sep) {
				return sep
			}
		}
	}
	return ""
}

//将YYYYMMDD格式的日期转换为YYYY-MM-DD格式或者YYYY-MM格式
func (h *TimeHelper) AddDateSeparator(date, separator string, incDay bool) string {
	temp := h.RemoveDateSeparator(date)
	if len(temp) < 8 {
		return date
	}

	if incDay {
		return fmt.Sprintf("%s%s%s%s%s", temp[:4], separator, temp[4:6], separator, temp[6:8])
	} else {
		return fmt.Sprintf("%s%s%s", temp[:4], separator, temp[4:6])
	}
}

//将日期时间字符串中的时间信息去掉
func (h *TimeHelper) RemoveTimeFromDateTime(date string) string {
	separator := h.GetDateSeparator(date)
	return h.AddDateSeparator(date, separator, true)
}

//获取短日期字符串
func (h *TimeHelper) GetShortDate(date string) string {
	temp := h.RemoveDateSeparator(date)
	if len(temp) < 8 {
		return date
	}
	return fmt.Sprintf("%s%s", temp[4:6], temp[6:8])
}

/**
* 获取当前日期时间字符串
 */
func (h *TimeHelper) GetNowString() string {
	return h.GetNowStringEx("yyyy-MM-dd HH:mm:ss", false)
}

/**
* 获取当前日期时间字符串
 */
func (h *TimeHelper) GetNowStringEx(format string, isUtcTime bool) string {
	t := time.Now()
	if isUtcTime {
		t = t.UTC()
	}
	return CcStr.FormatTime(t, format)
}

/**
* 将日期时间格式化字符串后再转换成64位整数
 */
func (h *TimeHelper) TimeToInt64(t time.Time, format string) int64 {
	s := CcStr.FormatTime(t, format)
	return CcStr.ParseInt64(s)
}

/**
* 将日期时间转换成 yyyy-MM-dd 00:00:00:000
 */
func (h *TimeHelper) ToDayStart(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

/**
* 将日期时间转换成 yyyy-MM-dd 23:59:59:999
 */
func (h *TimeHelper) ToDayEnd(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 59, t.Location())
}
