package ftime

import "go.mongodb.org/mongo-driver/bson/bsontype"

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	data, err := t.MarshalJSON()

	switch UsedFormatType {
	case TIMESTAMP:
		return bsontype.Int64, data, err
	case TEXT:
		return bsontype.String, data, err
	default:
		return bsontype.Int64, nil, ErrUnKnownFormatType
	}
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	switch bt {
	case bsontype.Int64:
		return t.parseTS(string(data))
	case bsontype.String:
		return t.parseText(string(data))
	default:
		return ErrUnKnownFormatType
	}
}
