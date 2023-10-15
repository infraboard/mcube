package ftime

import (
	"strconv"
	"time"
)

// FormatType 时间输出的格式
type FormatType string

const (
	// TIMESTAMP 时间戳展示
	TIMESTAMP = FormatType("timestamp")
	// TEXT 文本格式展示
	TEXT = FormatType("text")
)

// TimestampLength 时间戳长度
type TimestampLength int

const (
	// Length10 10位 精确到秒
	Length10 = TimestampLength(10)
	// Length13 13位 精确到毫秒
	Length13 = TimestampLength(13)
	// Length16 16位 精确到微妙
	Length16 = TimestampLength(16)
	// Length19 19位 精确到纳秒
	Length19 = TimestampLength(19)
)

var (
	// UsedFormatType 当前时间展示的格式
	UsedFormatType = TIMESTAMP
	// UsedTimestampLength 当前使用的时间戳长度
	UsedTimestampLength = Length13
	// UsedTextTimeFormart 文本格式
	UsedTextTimeFormart = "2006-01-02 15:04:05"
)

// Now 当前时间
func Now() Time {
	return T(time.Now())
}

// T 类型转换
func T(t time.Time) Time {
	return Time(t)
}

// Time 封装Time
// 主要用于序列化时控制输出格式
type Time time.Time

// T 转化成标准时间对象
func (t *Time) T() time.Time {
	return time.Time(*t)
}

func (t *Time) parseTS(ts string) error {
	// 补全到19位
	for len(ts) < 19 {
		ts += "0"
	}

	// 获取秒
	sec, err := strconv.ParseInt(ts[:10], 10, 64)
	if err != nil {
		return err
	}

	// 获取纳秒
	nsec, err := strconv.ParseInt(ts[10:], 10, 64)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(sec, nsec))
	return nil
}

func (t *Time) parseTSInt64(ts int64) error {
	return t.parseTS(strconv.FormatInt(ts, 10))
}

func (t *Time) parseText(data string) error {
	now, err := time.ParseInLocation(`"`+UsedTextTimeFormart+`"`, data, time.Local)
	if err != nil {
		return err
	}
	*t = Time(now)
	return nil
}

func (t Time) formatText() []byte {
	b := make([]byte, 0, len(UsedTextTimeFormart)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, UsedTextTimeFormart)
	b = append(b, '"')
	return b
}

// Timestamp 时间戳
func (t Time) Timestamp() int64 {
	var ts int64

	if t.T().IsZero() {
		return 0
	}

	switch UsedTimestampLength {
	case Length10:
		ts = time.Time(t).Unix()
	case Length13:
		ts = time.Time(t).UnixNano() / 1000000
	case Length16:
		ts = time.Time(t).UnixNano() / 1000
	case Length19:
		ts = time.Time(t).UnixNano()
	default:
		ts = time.Time(t).UnixNano() / 1000000
	}

	return ts
}
