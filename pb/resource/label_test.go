package resource_test

import (
	"testing"

	"github.com/infraboard/mcube/pb/resource"
)

func TestParseMapFromString(t *testing.T) {
	t.Log(resource.ParseMapFromString("key=value1,key2=value2=value3"))
}

func TestParseLabelRequirementFromString(t *testing.T) {
	t.Log(resource.ParseLabelRequirementFromString("key=value1,value2,value3"))
}

func TestParseLabelRequirementListFromString(t *testing.T) {
	t.Log(resource.ParseLabelRequirementListFromString("x1=y1,y2,y3&m1=n1,n2,n3"))
}
