package enum

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"text/template"

	"github.com/pkg/errors"
)

// G Generater
var G = Generater{
	t: template.New("enum"),
}

// Generater 用于生成枚举的生成器
type Generater struct {
	t           *template.Template
	file        string
	Marshal     bool
	ProtobufExt bool
}

// NewRenderParams todo
func NewRenderParams() *RenderParams {
	return &RenderParams{
		Enums:     NewEnumSet(),
		Stringer:  true,
		ValueMap:  true,
		Backquote: "`",
	}
}

// RenderParams 模板渲染需要的参数
type RenderParams struct {
	PKG       string
	Backquote string
	Enums     *Set
	Marshal   bool
	Stringer  bool
	ValueMap  bool
}

// Generate 生成文件
func (g *Generater) Generate(file string) ([]byte, error) {
	g.file = file
	params, err := g.parse()
	if err != nil {
		return nil, err
	}

	params.Marshal = g.Marshal
	if g.ProtobufExt {
		params.Stringer = false
		params.ValueMap = false
	}

	if params.Enums.Length() == 0 {
		return []byte{}, nil
	}

	return g.gen(params)
}

// 兼容命令行测试
func (g *Generater) getFile() string {
	if g.file != "" {
		return g.file
	}

	return os.Getenv("GOFILE")
}

// 解析代码源文件，获取常量和类型
func (g *Generater) parse() (*RenderParams, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, g.getFile(), nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parse file error, %s", err)
	}

	params := NewRenderParams()
	params.PKG = f.Name.Name

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			switch d.Tok {
			case token.CONST:
				for _, spec := range d.Specs {
					vs, _ := spec.(*ast.ValueSpec)

					ident := vs.Names[0]
					doc := vs.Doc.Text()

					var enum *Enum
					vst, _ := vs.Type.(*ast.Ident)

					if vst != nil {
						enum = params.Enums.Get(vst.Name)
						enum.Add(NewItem(ident.Name, doc))
					}

				}
			}
		}
	}

	return params, nil
}

func (g *Generater) gen(params *RenderParams) ([]byte, error) {
	buf := bytes.NewBufferString("")
	t, err := g.t.Parse(tmp)
	if err != nil {
		return nil, errors.Wrapf(err, "template init err")
	}

	err = t.Execute(buf, params)
	if err != nil {
		return nil, errors.Wrapf(err, "template data err")
	}
	return format.Source(buf.Bytes())
}
