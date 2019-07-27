package string

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

var (
	CcStr = NewStringHelper()
)

//string帮助类
type StringHelper struct {
}

//获取一个新的StringHelper实例
func NewStringHelper() *StringHelper {
	return &StringHelper{}
}

//格式化输入bool
func (h *StringHelper) FormatBool(b bool) string {
	return strconv.FormatBool(b)
}

//解析bool值
func (h *StringHelper) ParseBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

//格式化输出int值
func (h *StringHelper) FormatInt(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

//格式化输出int32值
func (h *StringHelper) FormatInt32(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

//格式化输出int64值
func (h *StringHelper) FormatInt64(i int64) string {
	return strconv.FormatInt(i, 10)
}

//解析int值
func (h *StringHelper) ParseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int(i)
}

//解析int32值
func (h *StringHelper) ParseInt32(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

//解析int64值
func (h *StringHelper) ParseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

//格式化输出uint值
func (h *StringHelper) FormatUint(u uint) string {
	return strconv.FormatUint(uint64(u), 10)
}

//格式化输出uint32值
func (h *StringHelper) FormatUint32(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

//格式化输出uint64值
func (h *StringHelper) FormatUint64(u uint64) string {
	return strconv.FormatUint(u, 10)
}

//解析uint值
func (h *StringHelper) ParseUint(s string) uint {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint(u)
}

//解析uint32值
func (h *StringHelper) ParseUint32(s string) uint32 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return uint32(u)
}

//解析uint64值
func (h *StringHelper) ParseUint64(s string) uint64 {
	u, _ := strconv.ParseUint(s, 10, 64)
	return u
}

//格式化输出float值
func (h *StringHelper) FormatFloat(f float32) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 32)

}

//格式化输出float64值
func (h *StringHelper) FormatFloat64(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

//解析float值
func (h *StringHelper) ParseFloat(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

//解析float64值
func (h *StringHelper) ParseFloat64(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

//格式化输出DateTime值，格式为：yyyy-MM-dd HH:mm:ss
func (h *StringHelper) FormatTime(t time.Time, format string) string {
	layout := h.TimeFormatToLayout(format)
	return t.Format(layout)
}

//解析DateTime值，格式为：yyyy-MM-dd HH:mm:ss
func (h *StringHelper) ParseTime(s, format string, utc bool) time.Time {
	layout := h.TimeFormatToLayout(format)
	if utc {
		t, _ := time.Parse(layout, s)
		return t
	} else {
		t, _ := time.ParseInLocation(layout, s, time.Now().Location())
		return t
	}
}

/**
* 从字符串中解析，解析失败时返回默认值
 */
func (h *StringHelper) ParseTimeEx(s, format string, utc bool, defaultTime time.Time) time.Time {
	t := h.ParseTime(s, format, utc)
	if t.IsZero() {
		return defaultTime
	}

	return t
}

/**
* 转换成字符串
 */
func (h *StringHelper) ToString(i interface{}) string {
	switch i.(type) {
	case int:
		return h.FormatInt(i.(int))
	case uint:
		return h.FormatUint(i.(uint))
	case int32:
		return h.FormatInt32(i.(int32))
	case uint32:
		return h.FormatUint32(i.(uint32))
	case int64:
		return h.FormatInt64(i.(int64))
	case uint64:
		return h.FormatUint64(i.(uint64))
	case float32:
		return h.FormatFloat(i.(float32))
	case float64:
		return h.FormatFloat64(i.(float64))
	case bool:
		return h.FormatBool(i.(bool))
	case time.Time:
		return h.FormatTime(i.(time.Time), "yyyy-MM-dd HH:mm:ss")
	case string:
		return i.(string)
	default:
		return "unknown type"
	}
}

//转换日期时间格式
func (h *StringHelper) ChangeDateTimeFormat(dateTime, srcFormat, dstFormat string) string {
	dt := h.ParseTime(dateTime, srcFormat, true)
	if dt.IsZero() {
		return dateTime
	} else {
		return h.FormatTime(dt, dstFormat)
	}
}

//将格式为：yyyy-MM-dd HH:mm:ss 转换 "2006-01-02 15:04:05"
func (h *StringHelper) TimeFormatToLayout(format string) string {
	layout := format
	layout = strings.Replace(layout, "yyyy", "2006", -1)
	layout = strings.Replace(layout, "yy", "06", -1)
	layout = strings.Replace(layout, "MMMM", "January", -1)
	layout = strings.Replace(layout, "MMM", "Jan", -1)
	layout = strings.Replace(layout, "MM", "01", -1)
	layout = strings.Replace(layout, "M", "1", -1)
	layout = strings.Replace(layout, "dddd", "Monday", -1)
	layout = strings.Replace(layout, "ddd", "Mon", -1)
	layout = strings.Replace(layout, "dd", "02", -1)
	layout = strings.Replace(layout, "d", "2", -1)
	layout = strings.Replace(layout, "HH", "15", -1)
	layout = strings.Replace(layout, "H", "15", -1)
	layout = strings.Replace(layout, "hh", "03", -1)
	layout = strings.Replace(layout, "h", "3", -1)
	layout = strings.Replace(layout, "mm", "04", -1)
	layout = strings.Replace(layout, "m", "4", -1)
	layout = strings.Replace(layout, "ss", "05", -1)
	layout = strings.Replace(layout, "s", "5", -1)
	layout = strings.Replace(layout, "fff", "999", -1)
	layout = strings.Replace(layout, "ff", "99", -1)
	layout = strings.Replace(layout, "f", "9", -1)
	layout = strings.Replace(layout, "zzz", "-0700", -1)
	layout = strings.Replace(layout, "zz", "-07", -1)
	layout = strings.Replace(layout, "z", "-07", -1)
	//layout = strings.Replace(layout, "ss", "0700", -1)
	return layout
}

//字符串去掉空格
func (h *StringHelper) TrimSpace(text string) string {
	if "" == text || len(text) < 0 {
		return ""
	}

	return strings.TrimSpace(text)
}

/**
* 字符串替换
 */
func (h *StringHelper) ReplaceAll(s, new string, old ...string) string {
	if "" == s || nil == old || len(old) <= 0 {
		return s
	}

	str := s

	for _, o := range old {
		str = strings.ReplaceAll(str, o, new)
	}

	return str
}

//获取子字符串
func (h *StringHelper) SubString(str string, begin, length int) (substr string) {
	// 将字符串的转换成[]rune
	rs := []rune(str)
	lth := len(rs)

	// 简单的越界判断
	if begin < 0 {
		begin = 0
	}

	if begin >= lth {
		begin = lth
	}

	end := begin + length
	if end > lth {
		end = lth
	}

	// 返回子串
	return string(rs[begin:end])
}

//获取Guid
func (h *StringHelper) NewGuid() (string, error) {
	guid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, guid)
	if n != len(guid) || err != nil {
		return "", err
	}

	guid[8] = guid[8]&^0xc0 | 0x80
	guid[6] = guid[6]&^0xf0 | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", guid[0:4], guid[4:6], guid[6:8], guid[8:10], guid[10:]), nil
}

//获取第一个有效的字符串
func (h *StringHelper) FirstValid(strArray ...string) string {
	if nil == strArray || len(strArray) <= 0 {
		return ""
	}

	var str string

	for _, str = range strArray {
		if str != "" {
			break
		}
	}

	return str
}
