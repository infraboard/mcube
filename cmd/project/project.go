package project

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"

	"github.com/infraboard/mcube/cmd/templates/api"
	"github.com/infraboard/mcube/cmd/templates/build"
	"github.com/infraboard/mcube/cmd/templates/cmd"
	"github.com/infraboard/mcube/cmd/templates/conf"
	"github.com/infraboard/mcube/cmd/templates/etc"
	"github.com/infraboard/mcube/cmd/templates/pkg"
	"github.com/infraboard/mcube/cmd/templates/root"
	"github.com/infraboard/mcube/cmd/templates/version"
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
	if err := p.rendTemplate("api", "api.go", api.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("version", "version.go", version.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("build", "build.sh", build.Template); err != nil {
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

	if err := p.rendTemplate("pkg", "auther.go", pkg.AutherTemplate); err != nil {
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

	return nil
}

func (p *Project) initGOModule() error {
	cmd := exec.Command("go", "mod", "init", p.PKG)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Printf("combined out:\n%s\n", string(out))
	return nil
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

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	t, err := p.render.Parse(tmpl)
	if err != nil {
		return fmt.Errorf("render %s/%s error, %s", dir, file, err)
	}
	return t.Execute(f, p)
}
