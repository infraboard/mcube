package project

import (
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func (p *Project) LoadMcenterConfig() error {
	p.Mcenter = &Mcenter{}

	var keyauthAddr string
	err := survey.AskOne(
		&survey.Input{
			Message: "Mcenter GRPC服务地址:",
			Default: "127.0.0.1:18050",
		},
		&keyauthAddr,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}
	if strings.Contains(keyauthAddr, ":") {
		hp := strings.Split(keyauthAddr, ":")
		p.Mcenter.Host = hp[0]
		p.Mcenter.Port = hp[1]
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "Keyauth Client ID:",
			Default: "",
		},
		&p.Mcenter.ClientID,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Password{
			Message: "Keyauth Client Secret:",
		},
		&p.Mcenter.ClientSecret,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) LoadMySQLConfig() error {
	p.MySQL = &MySQL{}

	var mySQLAddr string
	err := survey.AskOne(
		&survey.Input{
			Message: "MySQL服务地址:",
			Default: "127.0.0.1:3306",
		},
		&mySQLAddr,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	if strings.Contains(mySQLAddr, ":") {
		hp := strings.Split(mySQLAddr, ":")
		p.MySQL.Host = hp[0]
		p.MySQL.Port = hp[1]
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "数据库名称:",
			Default: "",
		},
		&p.MySQL.Database,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "用户:",
			Default: "",
		},
		&p.MySQL.UserName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Password{
			Message: "密码:",
		},
		&p.MySQL.Password,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) LoadMongoDBConfig() error {
	p.MongoDB = &MongoDB{}

	eps := ""
	err := survey.AskOne(
		&survey.Input{
			Message: "MongoDB服务地址,多个地址使用逗号分隔:",
			Default: "127.0.0.1:27017",
		},
		&eps,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}
	p.MongoDB.Endpoints = strings.Split(eps, ",")

	err = survey.AskOne(
		&survey.Input{
			Message: "认证数据库名称:",
			Default: "",
		},
		&p.MongoDB.AuthDB,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "认证用户:",
			Default: "",
		},
		&p.MongoDB.UserName,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Password{
			Message: "认证密码:",
		},
		&p.MongoDB.Password,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "数据库名称:",
			Default: p.MongoDB.AuthDB,
		},
		&p.MongoDB.Database,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return err
	}

	return nil
}
