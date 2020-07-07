package label_test

import (
	"testing"

	"github.com/infraboard/mcube/http/label"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	should := assert.New(t)

	l := label.Action("test")
	should.Equal(label.ActionLableKey, l.Key())
	should.Equal("test", l.Value())
}
