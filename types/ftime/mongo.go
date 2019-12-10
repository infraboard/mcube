package ftime

import "go.mongodb.org/mongo-driver/bson/bsontype"

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (t Time) MarshalBSONValue() (bsontype.Type, []byte, error) {
	data, err := t.MarshalJSON()
	return bsontype.Int64, data, err
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (t *Time) UnmarshalBSONValue(bt bsontype.Type, data []byte) error {
	return t.UnmarshalJSON(data)
}
