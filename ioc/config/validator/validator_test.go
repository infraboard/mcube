package validator_test

import (
	"testing"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/infraboard/mcube/v2/ioc/config/validator"
)

var (
	testStruct = Test{}
)

type Test struct {
	RequiredString   string   `validate:"required"`
	RequiredNumber   int      `validate:"required"`
	RequiredMultiple []string `validate:"required"`
}

func TestValidator(t *testing.T) {
	if err := validator.Validate(testStruct); err != nil {
		t.Fatal(err)
	}
}

func init() {
	err := ioc.ConfigIocObject(ioc.NewLoadConfigRequest())
	if err != nil {
		panic(err)
	}
}
