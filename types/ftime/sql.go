package ftime

import (
	"database/sql/driver"
)

// Scan 实现sql的反序列化
func (t Time) Scan(value interface{}) {
	switch UsedFormatType {
	case TIMESTAMP:
		value = t.timestamp()
	case TEXT:
		value = t.formatText()
	default:
		value = t.timestamp()
	}
}

// Value 实现sql的序列化
func (t *Time) Value() (driver.Value, error) {
	switch UsedFormatType {
	case TIMESTAMP:
		return t.timestamp(), nil
	case TEXT:
		return t.formatText(), nil
	default:
		return t.timestamp(), nil
	}

}
