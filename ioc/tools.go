package ioc

import (
	"fmt"
)

func ValidateIocObject(obj Object) error {
	if obj.Name() == "" {
		return fmt.Errorf("%T object name required", obj)
	}
	if obj.Version() == "" {
		return fmt.Errorf("%T object version required", obj)
	}
	return nil
}
