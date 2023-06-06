package ioc

import "fmt"

// 初始化托管的所有对象
func InitIocObject() error {
	for ns, objects := range store.store {
		objects.Sort()
		for i := range objects.Items {
			obj := objects.Items[i]
			err := obj.Init()
			if err != nil {
				return fmt.Errorf("init object %s.%s error, %s", ns, obj.Name(), err)
			}
		}
	}
	return nil
}
