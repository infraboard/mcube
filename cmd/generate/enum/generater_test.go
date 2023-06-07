package enum_test

import (
	"testing"

	"github.com/infraboard/mcube/cmd/generate/enum"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	should := assert.New(t)
	code, err := enum.G.Generate("../../../examples/enum/enum.go")
	t.Log(string(code))
	should.NoError(err)
}
