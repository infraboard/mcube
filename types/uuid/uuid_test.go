package uuid_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/infraboard/mcube/v2/types/uuid"
)

type TestStruct struct {
	ID uuid.BinaryUUID
}

func (s *TestStruct) String() string {
	return pretty.ToJSON(s)
}

func TestToJSON(t *testing.T) {
	s := TestStruct{
		ID: uuid.NewUUID(),
	}
	t.Log(s.String())
}
