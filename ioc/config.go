package ioc

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
	"github.com/infraboard/mcube/tools/file"
)

const (
	configNamespace = "configs"
)

// 用于托管配置对象的Ioc空间, 最先初始化
func Config() ObjectStroe {
	return store.
		Namespace(configNamespace).
		SetPriority(99)
}

// LoadConfig加载配置
func LoadConfig(req *LoadConfigRequest) error {
	objects := store.Namespace(configNamespace)
	errs := []string{}
	objects.ForEach(func(o Object) {
		var err error
		switch req.ConfigType {
		case CONFIG_TYPE_ENV:
			err = env.Parse(o, env.Options{
				Prefix: req.ConfigEnv.Prefix,
			})
		case CONFIG_TYPE_FILE:
			cfg := map[string]Object{
				o.Name(): o,
			}
			switch req.ConfigFile.ConfigFileFormat {
			case CONFIG_FILE_FORMAT_TOML:
				_, err = toml.DecodeFile(req.ConfigFile.Path, cfg)
			case CONFIG_FILE_FORMAT_YAML:
				err = file.ReadYamlFile(req.ConfigFile.Path, cfg)
			case CONFIG_FILE_FORMAT_JSON:
				err = file.ReadJsonFile(req.ConfigFile.Path, cfg)
			}
		}
		if err != nil {
			errs = append(errs, err.Error())
		}
	})

	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, ","))
	}
	return nil
}

func NewLoadConfigRequest() *LoadConfigRequest {
	return &LoadConfigRequest{
		ConfigType: CONFIG_TYPE_ENV,
		ConfigFile: &ConfigFile{
			ConfigFileFormat: CONFIG_FILE_FORMAT_TOML,
			Path:             "etc/application.toml",
		},
		ConfigEnv: &ConfigEnv{},
	}
}

type LoadConfigRequest struct {
	// 配置方式
	ConfigType ConfigType
	// 文件配置方式
	ConfigFile *ConfigFile
	// 环境变量配置
	ConfigEnv *ConfigEnv
}

type ConfigFile struct {
	ConfigFileFormat ConfigFileFormat
	Path             string
}

type ConfigType int

const (
	CONFIG_TYPE_ENV  ConfigType = iota
	CONFIG_TYPE_FILE            = 1
)

type ConfigFileFormat int

const (
	CONFIG_FILE_FORMAT_TOML ConfigFileFormat = iota
	CONFIG_FILE_FORMAT_YAML
	CONFIG_FILE_FORMAT_JSON
)

type ConfigEnv struct {
	Prefix string
}
