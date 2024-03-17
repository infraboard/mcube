package validator

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "validator"
)

func Validate(target any) error {
	return Get().Validate(target)
}

func Get() *Config {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Config)
}
