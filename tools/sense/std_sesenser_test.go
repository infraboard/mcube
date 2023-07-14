package sense_test

import (
	"testing"

	"github.com/infraboard/mcube/tools/sense"
)

func TestDeSense(t *testing.T) {
	t.Logf(sense.DeSense("abcdefghijkl"))
}
