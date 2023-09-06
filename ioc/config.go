package ioc

const (
	configNamespace = "configs"
)

// 用于托管配置对象的Ioc空间, 最先初始化
func Config() ObjectStroe {
	return store.
		Namespace(configNamespace).
		SetPriority(99)
}

// LoadConfigFromToml 从toml中添加配置文件, 并初始化全局对象
// func LoadConfigFromFile(filePath string) error {
// 	objects := store.Namespace(configNamespace)

// 	errs := []string{}
// 	objects.ForEach(func(o Object) {
// 		cfg := map[string]Object{
// 			o.Name(): o,
// 		}
// 		if _, err := toml.DecodeFile(filePath, cfg); err != nil {
// 			errs = append(errs, err.Error())
// 		}
// 	})

// 	if len(errs) > 0 {
// 		return fmt.Errorf("%s", strings.Join(errs, ","))
// 	}
// 	return nil
// }
