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
func Config() Stroe {
	return store.
		Namespace(configNamespace).
		SetPriority(99)
}

// LoadConfig加载配置
func LoadConfig(req *LoadConfigRequest) error {
	objects := store.Namespace(configNamespace)
	errs := []string{}
	objects.ForEach(func(o Object) {

		err := env.Parse(o, env.Options{
			Prefix: req.ConfigEnv.Prefix,
		})
		if err != nil {
			errs = append(errs, err.Error())
		}

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
		ConfigEnv: &ConfigEnv{},
		ConfigFile: &ConfigFile{
			ConfigFileFormat: CONFIG_FILE_FORMAT_TOML,
			Path:             "etc/application.toml",
		},
	}
}

type LoadConfigRequest struct {
	// 环境变量配置
	ConfigEnv *ConfigEnv
	// 文件配置方式
	ConfigFile *ConfigFile
}

type ConfigFile struct {
	Enabled          bool
	ConfigFileFormat ConfigFileFormat
	Path             string
}

type ConfigFileFormat int

const (
	CONFIG_FILE_FORMAT_TOML ConfigFileFormat = iota
	CONFIG_FILE_FORMAT_YAML
	CONFIG_FILE_FORMAT_JSON
)

type ConfigEnv struct {
	Enabled bool
	Prefix  string
}
