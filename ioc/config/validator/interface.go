package validator

import "github.com/infraboard/mcube/v2/ioc"

const (
	AppName = "validator"
)

func Validate(target any) error {
	return ioc.Config().Get(AppName).(*Config).Validate(target)
}
