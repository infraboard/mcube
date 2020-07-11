package ftime

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	var dst []byte
	switch UsedFormatType {
	case TIMESTAMP:
		dst = bsoncore.AppendInt64(dst, t.Timestamp())
		return bsontype.Int64, dst, nil
	case TEXT:
		dst = bsoncore.AppendString(dst, string(t.formatText()))
		return bsontype.String, dst, nil
	default:
		return bsontype.Int64, nil, ErrUnKnownFormatType
	}
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	if t == nil {
		return errors.New("cannot unmarshal into nil Value")
	}

	switch bt {
	case bsontype.Int64:
		var i64 int64
		i64, rem, ok := bsoncore.ReadInt64(data)
		if !ok {
			return bsoncore.NewInsufficientBytesError(data, rem)
		}
		return t.parseTSInt64(i64)
	case bsontype.String:
		str, rem, ok := bsoncore.ReadString(data)
		if !ok {
			return bsoncore.NewInsufficientBytesError(data, rem)
		}
		return t.parseText(string(str))
	default:
		return ErrUnKnownFormatType
	}

}
