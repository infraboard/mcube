package ioc

import "reflect"

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
