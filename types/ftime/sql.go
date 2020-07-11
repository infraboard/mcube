package ftime

import (
	"database/sql/driver"
	"errors"
)

// Scan 实现sql的反序列化
func (t *Time) Scan(value interface{}) error {
	if t == nil {
		return nil
	}

	switch v := value.(type) {
	case int64:
		return t.parseTSInt64(v)
	case string:
		return t.parseTS(v)
	default:
		return errors.New("unsupport type")
	}
}

// Value 实现sql的序列化
func (t Time) Value() (driver.Value, error) {
	switch UsedFormatType {
	case TIMESTAMP:
		return t.Timestamp(), nil
	case TEXT:
		return t.formatText(), nil
	default:
		return t.Timestamp(), nil
	}

}
