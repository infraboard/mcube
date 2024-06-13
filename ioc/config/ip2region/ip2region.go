package ip2region

import (
	"fmt"

	"github.com/infraboard/mcube/v2/ioc"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func init() {
	ioc.Config().Registry(defaultConfig)
}

var defaultConfig = &Ip2Region{
	Enable: true,
	DBPath: "etc/ip2region.xdb",
}

type Ip2Region struct {
	ioc.ObjectImpl

	// 功能开关, 开启后 需要读取DB文件, 在执行单元测试时很不方便
	Enable bool `json:"enable" yaml:"enable" toml:"enable" env:"ENABLE"`
	// DB 文件路径
	DBPath string `json:"db_path" yaml:"db_path" toml:"db_path" env:"DB_PATH"`

	searcher *xdb.Searcher
}

func (i *Ip2Region) Name() string {
	return AppName
}

func (i *Ip2Region) Init() error {
	if !i.Enable {
		return nil
	}

	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := xdb.LoadContentFromFile(i.DBPath)
	if err != nil {
		return fmt.Errorf("failed to load content from `%s`: %s", i.DBPath, err)
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	searcher, err := xdb.NewWithBuffer(cBuff)
	if err != nil {
		return fmt.Errorf("failed to create searcher with content: %s", err)
	}
	i.searcher = searcher
	return nil
}

func (i *Ip2Region) LookupIP(ip string) (*IPInfo, error) {
	if !i.Enable {
		return nil, fmt.Errorf("not enabled")
	}

	if i.searcher == nil {
		return nil, fmt.Errorf("ip lookup searcher is nil")
	}

	resp, err := i.searcher.SearchByStr(ip)
	if err != nil {
		return nil, err
	}
	return ParseIpInfoFromString(resp)
}
