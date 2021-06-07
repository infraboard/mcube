package main

import (
	"flag"

	"github.com/infraboard/mcube/cmd/protoc-gen-go-ext/ast"
	gengo "google.golang.org/protobuf/cmd/protoc-gen-go/internal_gengo"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var (
		flags        flag.FlagSet
		importPrefix = flags.String("import_prefix", "", "prefix to prepend to import paths")
	)
	importRewriteFunc := func(importPath protogen.GoImportPath) protogen.GoImportPath {
		switch importPath {
		case "context", "fmt", "math":
			return importPath
		}
		if *importPrefix != "" {
			return protogen.GoImportPath(*importPrefix) + importPath
		}
		return importPath
	}

	protogen.Options{
		ParamFunc:         flags.Set,
		ImportRewriteFunc: importRewriteFunc,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = gengo.SupportedFeatures
		var originFiles []*protogen.GeneratedFile
		for _, f := range gen.Files {
			if f.Generate {
				originFiles = append(originFiles, gengo.GenerateFile(gen, f))
			}
		}

		ast.Rewrite(gen)
		for _, f := range originFiles {
			f.Skip()
		}
		return nil
	})
}
