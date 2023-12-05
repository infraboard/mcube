package validator

import (
	"fmt"
	"strings"

	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_trans "github.com/go-playground/validator/v10/translations/zh"
	"github.com/infraboard/mcube/v2/ioc"
)

func init() {
	ioc.Config().Registry(&Config{})
}

type Config struct {
	ioc.ObjectImpl

	v     *validator.Validate
	trans ut.Translator
}

func (m *Config) Name() string {
	return AppName
}

func (m *Config) Init() error {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	tr, ok := uni.GetTranslator("zh")
	if !ok {
		return fmt.Errorf("zh not found")
	}

	m.trans = tr
	validate := validator.New()

	err := zh_trans.RegisterDefaultTranslations(validate, m.trans)
	if err != nil {
		return err
	}

	m.v = validate
	return nil
}

func (m *Config) Validate(target any) error {
	if m.v == nil {
		return fmt.Errorf("validator not init")
	}

	err := m.v.Struct(target)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fe := errs.Translate(m.trans)
	errStr := []string{}
	for _, v := range fe {
		errStr = append(errStr, v)
	}
	if len(errStr) > 0 {
		return fmt.Errorf(strings.Join(errStr, ","))
	}

	return nil
}
