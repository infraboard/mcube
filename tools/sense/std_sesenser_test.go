package sense_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/tools/sense"
)

func TestDeSense(t *testing.T) {
	t.Logf(sense.DeSense("abcdefghijkl"))
}
