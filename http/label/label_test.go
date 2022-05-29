package label_test

import (
	"testing"

	"github.com/infraboard/mcube/http/label"
	"github.com/stretchr/testify/assert"
)

func TestAction(t *testing.T) {
	should := assert.New(t)

	l := label.NewActionLabel("test")
	should.Equal(label.Action, l.Key())
	should.Equal("test", l.Value())
}
