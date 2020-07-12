package router_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/infraboard/mcube/http/label"
	"github.com/infraboard/mcube/http/router"
)

func TestEntry(t *testing.T) {
	should := assert.New(t)

	e := router.NewEntry("/mcube/v1/", "GET", "Monkey")
	e.EnableAuth()
	e.EnablePermission()
	e.AddLabel(label.Get)

	should.Equal("/mcube/v1/ [GET] Monkey map[action:get]", e.String())

	set := router.NewEntrySet()
	set.AddEntry(*e, *e)
	should.Equal("/mcube/v1/ [GET] Monkey map[action:get]\n/mcube/v1/ [GET] Monkey map[action:get]", set.String())
}
