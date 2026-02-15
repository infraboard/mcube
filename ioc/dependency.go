package ioc

import (
	"fmt"
	"reflect"
	"strings"
)

// PrintDependencies 打印命名空间的依赖关系（树形结构）
func (s *NamespaceStore) PrintDependencies() {
	fmt.Printf("\n=== %s Namespace Dependency Tree ===\n", s.Namespace)

	visited := make(map[string]bool)

	s.ForEach(func(w *ObjectWrapper) {
		if !visited[w.Name] {
			fmt.Print("├─ ")
			s.printDependencyTree(w.Name, w.Value, "   ", visited, make(map[string]bool))
		}
	})

	fmt.Println()
}

// printDependencyTree 递归打印依赖树
func (s *NamespaceStore) printDependencyTree(name string, obj Object, prefix string, visited, inPath map[string]bool) {
	// 检测循环依赖
	if inPath[name] {
		fmt.Printf("%s@%s ⚠️  (circular dependency)\n", name, obj.Version())
		return
	}

	// 打印当前对象（带版本号）
	version := obj.Version()
	deps := s.extractDependencies(obj)

	if len(deps) == 0 {
		fmt.Printf("%s@%s\n", name, version)
		visited[name] = true
		return
	}

	fmt.Printf("%s@%s (%d deps)\n", name, version, len(deps))
	visited[name] = true
	inPath[name] = true

	// 递归打印依赖
	for i, depInfo := range deps {
		isLast := i == len(deps)-1
		connector := "├─"
		childPrefix := prefix + "│  "

		if isLast {
			connector = "└─"
			childPrefix = prefix + "   "
		}

		fmt.Printf("%s%s ", prefix, connector)

		// 查找依赖对象
		depObj := s.findDependencyObject(depInfo)
		if depObj != nil {
			depVersion := depObj.Version()

			if !visited[depInfo.Name] {
				// 递归打印未访问过的依赖
				s.printDependencyTree(depInfo.Name, depObj, childPrefix, visited, inPath)
			} else {
				// 已访问过，只显示引用
				fmt.Printf("%s@%s (already shown)\n", depInfo.Name, depVersion)
			}
		} else {
			fmt.Printf("%s (not found)\n", depInfo.Name)
		}
	}

	delete(inPath, name)
}

// DependencyInfo 依赖信息
type DependencyInfo struct {
	Name      string // 依赖对象名称
	Namespace string // 依赖所在命名空间
	FieldName string // 字段名
}

// extractDependencies 提取对象的依赖信息
// 支持两种方式：
// 1. 声明式依赖：通过 ioc 标签自动检测（推荐）
// 2. 命令式依赖：实现 DependencyDeclarer 接口手动声明
func (s *NamespaceStore) extractDependencies(obj Object) []DependencyInfo {
	var deps []DependencyInfo

	// 方式1：检查是否实现了 DependencyDeclarer 接口（命令式依赖声明）
	if declarer, ok := obj.(DependencyDeclarer); ok {
		declared := declarer.DeclareDependencies()
		deps = append(deps, declared...)
	}

	// 方式2：扫描 ioc 标签（声明式依赖，自动检测）
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("ioc")

		if tag == "" {
			continue
		}

		// 解析ioc标签
		tagInfo, err := ParseInjectTagWithError(tag)
		if err != nil {
			continue
		}

		if tagInfo.Autowire {
			depInfo := DependencyInfo{
				FieldName: field.Name,
				Namespace: tagInfo.Namespace,
			}

			// 确定依赖对象的名称
			if tagInfo.Name != "" {
				depInfo.Name = tagInfo.Name
			} else {
				depInfo.Name = field.Type.String()
			}

			deps = append(deps, depInfo)
		}
	}

	return deps
}

// findDependencyObject 查找依赖对象
func (s *NamespaceStore) findDependencyObject(depInfo DependencyInfo) Object {
	// 如果指定了namespace，从对应的namespace查找
	if depInfo.Namespace != "" && depInfo.Namespace != s.Namespace {
		targetNs := DefaultStore.Namespace(depInfo.Namespace)
		return targetNs.Get(depInfo.Name)
	}

	// 否则在当前namespace查找
	return s.Get(depInfo.Name)
}

// PrintAllDependencies 打印所有命名空间的依赖关系
func (d *defaultStore) PrintAllDependencies() {
	fmt.Println("\n╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║              IOC Container Dependency Tree                  ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")

	for _, ns := range d.store {
		ns.PrintDependencies()
	}
}

// PrintDependencySummary 打印依赖统计摘要
func (s *NamespaceStore) PrintDependencySummary() {
	fmt.Printf("\n=== %s Namespace Summary ===\n", s.Namespace)

	totalObjects := s.Len()
	depCounts := make(map[string]int)
	usedBy := make(map[string][]string)

	// 收集依赖信息
	s.ForEach(func(w *ObjectWrapper) {
		deps := s.extractDependencies(w.Value)
		depCounts[w.Name] = len(deps)

		for _, dep := range deps {
			usedBy[dep.Name] = append(usedBy[dep.Name], w.Name)
		}
	})

	// 统计
	noDeps := 0
	for _, count := range depCounts {
		if count == 0 {
			noDeps++
		}
	}

	fmt.Printf("  📊 Total Objects: %d\n", totalObjects)
	fmt.Printf("  🌿 Leaf Objects (no deps): %d\n", noDeps)
	fmt.Printf("  🔗 Objects with deps: %d\n", totalObjects-noDeps)

	// 找出依赖最多的对象
	maxDeps := 0
	var maxDepObj string
	for name, count := range depCounts {
		if count > maxDeps {
			maxDeps = count
			maxDepObj = name
		}
	}
	if maxDeps > 0 {
		fmt.Printf("  ⬆️  Most dependencies: %s (%d deps)\n", maxDepObj, maxDeps)
	}

	// 找出被使用最多的对象
	maxUsed := 0
	var maxUsedObj string
	for name, users := range usedBy {
		if len(users) > maxUsed {
			maxUsed = len(users)
			maxUsedObj = name
		}
	}
	if maxUsed > 0 {
		fmt.Printf("  ⬇️  Most depended on: %s (used by %d objects)\n", maxUsedObj, maxUsed)
	}

	fmt.Println()
}

// ExportDependenciesToMarkdown 导出依赖关系为Markdown格式
func (s *NamespaceStore) ExportDependenciesToMarkdown() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# %s Namespace Dependencies\n\n", s.Namespace))
	sb.WriteString("## Objects\n\n")

	s.ForEach(func(w *ObjectWrapper) {
		deps := s.extractDependencies(w.Value)

		sb.WriteString(fmt.Sprintf("### %s@%s\n\n", w.Name, w.Version))

		if len(deps) == 0 {
			sb.WriteString("- No dependencies\n\n")
		} else {
			sb.WriteString("**Dependencies:**\n\n")
			for _, dep := range deps {
				depObj := s.findDependencyObject(dep)
				if depObj != nil {
					sb.WriteString(fmt.Sprintf("- `%s@%s` (field: `%s`", dep.Name, depObj.Version(), dep.FieldName))
					if dep.Namespace != "" {
						sb.WriteString(fmt.Sprintf(", namespace: `%s`", dep.Namespace))
					}
					sb.WriteString(")\n")
				} else {
					sb.WriteString(fmt.Sprintf("- `%s` (field: `%s`, ⚠️ not found)\n", dep.Name, dep.FieldName))
				}
			}
			sb.WriteString("\n")
		}
	})

	return sb.String()
}
