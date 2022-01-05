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

	"github.com/infraboard/mcube/cmd/mcube/templates/app"
	"github.com/infraboard/mcube/cmd/mcube/templates/client"
	"github.com/infraboard/mcube/cmd/mcube/templates/cmd"
	"github.com/infraboard/mcube/cmd/mcube/templates/conf"
	"github.com/infraboard/mcube/cmd/mcube/templates/etc"
	"github.com/infraboard/mcube/cmd/mcube/templates/protocol"
	"github.com/infraboard/mcube/cmd/mcube/templates/root"
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

	// 选择是否接入权限中心Keyauth
	enableKeyauth := &survey.Confirm{
		Message: "是否接入权限中心[keyauth]",
	}
	survey.AskOne(enableKeyauth, &p.EnableKeyauth)

	if p.EnableKeyauth {
		p.LoadKeyauthConfig()
	}

	// 选择使用的存储
	choicedDB := ""
	choiceDB := &survey.Select{
		Message: "选择数据库类型:",
		Options: []string{"MySQL", "MongoDB"},
		Default: "MySQL",
	}
	survey.AskOne(choiceDB, &choicedDB)

	switch choicedDB {
	case "MySQL":
		p.LoadMySQLConfig()
	case "MongoDB":
		p.LoadMongoDBConfig()
	}

	p.caculate()
	return p, nil
}

// Project todo
type Project struct {
	PKG           string
	Name          string
	Description   string
	Backquote     string
	Backquote3    string
	EnableKeyauth bool
	Keyauth       *Keyauth
	EnableMySQL   bool
	MySQL         *MySQL
	EnableMongoDB bool
	MongoDB       *MongoDB

	render     *template.Template
	createdDir map[string]bool
}

// Keyauth 鉴权服务配置
type Keyauth struct {
	Host         string
	Port         string
	ClientID     string
	ClientSecret string
}

type MySQL struct {
	Host     string
	Port     string
	Database string
	UserName string
	Password string
}

type MongoDB struct {
	Endpoints string
	UserName  string
	Password  string
	Database  string
	AuthDB    string
}

func (p *Project) caculate() {
	if p.PKG != "" {
		slice := strings.Split(p.PKG, "/")
		p.Name = slice[len(slice)-1]
	}
}

// Init 初始化项目
func (p *Project) Init() error {
	if err := p.rendTemplate("protocol", "http.go", protocol.HTTPTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("protocol", "grpc.go", protocol.GRPCTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("version", "version.go", version.Template); err != nil {
		return err
	}

	if err := p.rendTemplate("cmd", "root.go", cmd.RootTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("cmd", "start.go", cmd.StartTemplate); err != nil {
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

	if err := p.rendTemplate("app/all", "all.go", app.AllTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/pb", "reponse.proto", app.ExamplePBResponseTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/pb", "request.proto", app.ExamplePBRequestTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/pb", "service.proto", app.ExamplePBServiceTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example", "request_ext.go", app.ExampleRequestExtTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example", "response_ext.go", app.ExampleResponseExtTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/impl", "impl.go", app.ExampleIMPLOBJTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/impl", "example.go", app.ExampleIMPLMethodTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/http", "http.go", app.ExampleHTTPObjTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("app/example/http", "example.go", app.ExampleHTTPMethodTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("client", "client.go", client.ClientProxyTemplate); err != nil {
		return err
	}

	if err := p.rendTemplate("client", "config.go", client.ClientConfigTemplate); err != nil {
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
			err := os.MkdirAll(dir, os.ModePerm)
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
