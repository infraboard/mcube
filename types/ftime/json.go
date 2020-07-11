package ftime

import (
	"encoding/json"
)

// MarshalJSON 实现JSON 序列化接口
func (t Time) MarshalJSON() ([]byte, error) {
	switch UsedFormatType {
	case TIMESTAMP:
		return json.Marshal(t.Timestamp())
	case TEXT:
		return json.Marshal(t.formatText())
	default:
		return nil, ErrUnKnownFormatType
	}
}

// UnmarshalJSON 实现JSON 反序列化接口
func (t *Time) UnmarshalJSON(b []byte) error {
	switch UsedFormatType {
	case TIMESTAMP:
		return t.parseTS(string(b))
	case TEXT:
		return t.parseText(string(b))
	default:
		return ErrUnKnownFormatType
	}
}
