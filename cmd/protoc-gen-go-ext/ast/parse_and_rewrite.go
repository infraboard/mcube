// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ast

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/searKing/golang/go/reflect"
	strings_ "github.com/searKing/golang/go/strings"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"

	pb "github.com/infraboard/mcube/cmd/protoc-gen-go-ext/extension/tag"
)

type FieldInfo struct {
	FieldNameInProto string
	FieldNameInGo    string
	FieldTag         reflect.StructTag
}
type StructInfo struct {
	StructNameInProto string
	StructNameInGo    string
	FieldInfos        []FieldInfo
}

type FileInfo struct {
	FileName    string
	StructInfos []StructInfo
	HasEnum     bool
}

func (si *StructInfo) FindField(name string) (FieldInfo, bool) {
	for _, f := range si.FieldInfos {
		if f.FieldNameInGo == name {
			return f, true
		}
	}
	return FieldInfo{}, false
}

func WalkDescriptorProto(g *protogen.Plugin, dp *descriptor.DescriptorProto, typeNames []string) []StructInfo {
	var ss []StructInfo

	s := StructInfo{}
	s.StructNameInProto = dp.GetName()
	s.StructNameInGo = CamelCaseSlice(append(typeNames, CamelCase(dp.GetName())))

	//typeNames := []string{s.StructNameInGo}
	for _, field := range dp.GetField() {
		if field.GetOptions() == nil {
			continue
		}

		v := proto.GetExtension(field.Options, pb.E_FieldTag)
		switch v := v.(type) {
		case *pb.FieldTag:
			tag := v.GetStructTag()
			tags, err := reflect.ParseStructTag(tag)
			if err != nil {
				g.Error(fmt.Errorf("failed to parse struct tag in field extension: %w", err))
				// ignore this tag
				continue
			}

			s.FieldInfos = append(s.FieldInfos, FieldInfo{
				FieldNameInProto: field.GetName(),
				FieldNameInGo:    CamelCase(field.GetName()),
				FieldTag:         *tags,
			})
		}
	}
	if len(s.FieldInfos) > 0 {
		ss = append(ss, s)
	}

	typeNames = append(typeNames, CamelCase(dp.GetName()))
	for _, nest := range dp.GetNestedType() {
		nestSs := WalkDescriptorProto(g, nest, typeNames)
		if len(nestSs) > 0 {
			ss = append(ss, nestSs...)
		}
	}
	return ss
}

// Rewrite 重新文件
func Rewrite(g *protogen.Plugin) {
	var protoFiles []FileInfo

	p := newProtoGenParams(g.Request.GetParameter())

	for _, protoFile := range g.Request.GetProtoFile() {
		if !strings_.SliceContains(g.Request.GetFileToGenerate(), protoFile.GetName()) {
			continue
		}
		f := FileInfo{}
		f.FileName = protoFile.GetName()

		if len(protoFile.EnumType) > 0 {
			f.HasEnum = true
		}

		for _, messageType := range protoFile.GetMessageType() {
			ss := WalkDescriptorProto(g, messageType, nil)
			if len(ss) > 0 {
				f.StructInfos = append(f.StructInfos, ss...)
			}
		}

		protoFiles = append(protoFiles, f)
	}
	// g.Response() will generate files, so skip this step
	//if len(g.Response().GetFile()) == 0 {
	//	return
	//}

	rewriter := NewGenerator(protoFiles, g)
	for _, f := range g.Response().GetFile() {
		rewriter.ParseGoContent(f)
	}
	rewriter.Generate(p.Get("module"))
}

func newProtoGenParams(s string) *parameter {
	kvmap := map[string]string{}
	kvStr := strings.Split(s, ",")
	for _, kv := range kvStr {
		index := strings.Index(kv, "=")
		if index > 0 {
			kvmap[kv[:index]] = kv[index+1:]
		}
	}
	return &parameter{kv: kvmap}
}

type parameter struct {
	kv map[string]string
}

func (p *parameter) Get(key string) string {
	if v, ok := p.kv[key]; ok {
		return v
	}

	return ""
}
