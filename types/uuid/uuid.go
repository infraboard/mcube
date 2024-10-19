package uuid

import (
	"database/sql/driver"

	"github.com/google/uuid"
)

// BinaryUUID -> binary uuid wrapper over uuid.UUID
type BinaryUUID uuid.UUID

// ParseUUID -> parses string uuid to binary uuid
func ParseUUID(id string) BinaryUUID {
	return BinaryUUID(uuid.MustParse(id))
}

func NewUUID() BinaryUUID {
	return BinaryUUID(uuid.Must(uuid.NewV7()))
}

func (b BinaryUUID) String() string {
	return uuid.UUID(b).String()
}

// MarshalJSON -> convert to json string
func (b BinaryUUID) MarshalJSON() ([]byte, error) {
	return uuid.UUID(b).MarshalText()
}

// UnmarshalJSON -> convert from json string
func (b *BinaryUUID) UnmarshalJSON(by []byte) error {
	s, err := uuid.ParseBytes(by)
	*b = BinaryUUID(s)
	return err
}

// GormDataType -> sql data type for gorm
func (BinaryUUID) GormDataType() string {
	return "binary(16)"
}

// Scan -> scan value into BinaryUUID
func (b *BinaryUUID) Scan(value any) error {
	s := uuid.Nil
	if err := s.Scan(value); err != nil {
		return err
	}
	*b = BinaryUUID(s)
	return nil
}

// Value -> return BinaryUUID to []bytes binary(16)
func (b BinaryUUID) Value() (driver.Value, error) {
	return uuid.UUID(b).MarshalBinary()
}
