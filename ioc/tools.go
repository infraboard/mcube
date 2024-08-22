package ioc

import (
	"os"
	"reflect"
)

func GetIocObjectUid(obj Object) (name, version string) {
	name = obj.Name()
	if name == "" {
		name = reflect.TypeOf(obj).String()
	}
	version = obj.Version()
	if version == "" {
		version = DEFAULT_VERSION
	}
	return
}

func IsFileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
