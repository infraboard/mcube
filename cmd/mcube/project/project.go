package project

import (
	"bytes"
	"fmt"
	"go/format"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/AlecAivazis/survey/v2"
	"github.com/pkg/errors"

	"github.com/infraboard/mcube/tools/cli"

	"embed"
)

//go:embed templates/*
var templates embed.FS

// LoadConfigFromCLI 配置
func LoadConfigFromCLI() (*Project, error) {
	p := &Project{
		render:     template.New("project"),
		createdDir: map[string]bool{},
	}

	err := survey.AskOne(
		&survey.Input{
			Message: "请输入项目包名称:",
			Default: "gitee.com/go-course/mcube-demo",
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
	err = survey.AskOne(enableKeyauth, &p.EnableKeyauth)
	if err != nil {
		return nil, err
	}

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
	err = survey.AskOne(choiceDB, &choicedDB)
	if err != nil {
		return nil, err
	}

	switch choicedDB {
	case "MySQL":
		p.EnableMySQL = true
		p.LoadMySQLConfig()
	case "MongoDB":
		p.EnableMongoDB = true
		p.LoadMongoDBConfig()
	}

	// 选择是否生成样例
	genExample := &survey.Confirm{
		Message: "生成样例代码",
	}
	survey.AskOne(genExample, &p.GenExample)

	p.caculate()
	return p, nil
}

// Project todo
type Project struct {
	PKG           string
	Name          string
	Description   string
	EnableKeyauth bool
	Keyauth       *Keyauth
	EnableMySQL   bool
	MySQL         *MySQL
	EnableMongoDB bool
	MongoDB       *MongoDB
	GenExample    bool
	EnableCache   bool

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
	fn := func(path string, d fs.DirEntry, _ error) error {
		// 不处理目录
		if d.IsDir() {
			return nil
		}

		// 忽略不是模板的文件
		if !strings.HasSuffix(d.Name(), ".tpl") {
			return nil
		}

		// 读取模板内容
		data, err := templates.ReadFile(path)
		if err != nil {
			return err
		}

		// 替换templates为项目目录名称
		target := strings.Replace(path, "templates", p.Name, 1)
		dirName := filepath.Dir(target)
		sourceFileName := strings.TrimSuffix(filepath.Base(target), ".tpl")

		return p.rendTemplate(dirName, sourceFileName, string(data))
	}

	err := fs.WalkDir(templates, "templates", fn)
	if err != nil {
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
	return cli.Tree(os.Stdout, p.Name, true)
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
