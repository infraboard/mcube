package enum_test

import (
	"testing"

	"github.com/infraboard/mcube/cmd/enum"
	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	should := assert.New(t)
	enum.G.SetSrcFile("./example_test.go")
	code, err := enum.G.Generate()
	t.Log(string(code))
	should.NoError(err)
}
