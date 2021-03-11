package project

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/infraboard/mcube/cmd/mcube/templates/api"
	"github.com/infraboard/mcube/cmd/mcube/templates/cmd"
	"github.com/infraboard/mcube/cmd/mcube/templates/conf"
	"github.com/infraboard/mcube/cmd/mcube/templates/etc"
	"github.com/infraboard/mcube/cmd/mcube/templates/pkg"
	"github.com/infraboard/mcube/cmd/mcube/templates/root"
	"github.com/infraboard/mcube/cmd/mcube/templates/script"
	"github.com/infraboard/mcube/cmd/mcube/templates/version"
	"github.com/infraboard/mcube/tools/cli"
)

// LoadConfigFromCLI 配置
func LoadConfigFromCLI() (*Project, error) {
	p := &Project{
		render:     template.New("project"),
		createdDir: map[string]bool{},
		Backquote:  "`",
		Backquote3: "```",
	}

	err := survey.AskOne(
		&survey.Input{
			Message: "请输入项目包名称:",
			Default: "github.com/infraboard/demo",
		},
		&p.PKG,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return nil, err
	}

	err = survey.AskOne(
		&survey.Input{
			Message: "请输入项目描述:",
			Default: "",
		},
		&p.Description,
		survey.WithValidator(survey.Required),
	)
	if err != nil {
		return nil, err
	}

	p.caculate()
	return p, nil
}

// Project todo
type Project struct {
	PKG         string
	Name        string
	Description string
	Backquote   string
	Backquote3  string

	render     *template.Template
	createdDir map[string]bool
}

func (p *Project) caculate() {
	if p.PKG != "" {
		slice := strings.Split(p.PKG, "/")
		p.Name = slice[len(slice)-1]
	}
}

// Init 初始化项目
func (p *Project) Init() error {
	if err := p.rendTemplate("api", "api.go", api.HTTPTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("api", "grpc.go", api.GRPCTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("version", "version.go", version.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("script", "build.sh", script.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("cmd", "root.go", cmd.RootTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("cmd", "service.go", cmd.ServiceTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("conf", "config.go", conf.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("conf", "load.go", conf.LoadTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("conf", "log.go", conf.LogTempate); err != nil {
		return err
	}

	if err := p.rendTemplate("pkg", "http.go", pkg.HTTPTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("pkg", "service.go", pkg.ServiceTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("etc", p.Name+".toml", etc.TOMLExampleTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("etc", p.Name+".env", etc.EnvExampleTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("", "main.go", root.MainTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("", "Makefile", root.MakefileTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("", "README.md", root.ReadmeTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("", ".gitignore", root.GitIgnreTemplate); err != nil {
		return err
	}

	if err := p.initGOModule(); err != nil {
		return err
	}

	fmt.Println("项目初始化完成, 项目结构如下: ")
	if err := p.show(); err != nil {
		return err
	}

	return nil
}

func (p *Project) initGOModule() error {
	cmd := exec.Command("go", "mod", "init", p.PKG)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) show() error {
	return cli.Tree(os.Stdout, ".", true)
}

func (p *Project) dirNotExist(path string) bool {
	if _, ok := p.createdDir[path]; ok {
		return false
	}

	return true
}

func (p *Project) rendTemplate(dir, file, tmpl string) error {
	if dir != "" {
		if p.dirNotExist(dir) {
			err := os.Mkdir(dir, os.ModePerm)
			if err != nil {
				return err
			}
			p.createdDir[dir] = true
		}
	}

	filePath := ""
	if dir != "" {
		filePath = dir + "/" + file
	} else {
		filePath = file
	}

	t, err := p.render.Parse(tmpl)
	if err != nil {
		return fmt.Errorf("render %s/%s error, %s", dir, file, err)
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, p)
	if err != nil {
		return errors.Wrapf(err, "template data err")
	}

	var content []byte
	if path.Ext(file) == "go" {
		code, err := format.Source(buf.Bytes())
		if err != nil {
			return errors.Wrapf(err, "format %s code err", file)
		}
		content = code
	} else {
		content = buf.Bytes()
	}

	return ioutil.WriteFile(filePath, content, 0644)
}
