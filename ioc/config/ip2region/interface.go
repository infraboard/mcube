package ip2region

import (
	"fmt"
	"strings"

	"github.com/infraboard/mcube/v2/ioc"
)

const (
	AppName = "ip2region"
)

func Get() IpRegionSearcher {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*Ip2Region)
}

type IpRegionSearcher interface {
	LookupIP(ip string) (*IPInfo, error)
}

// NewDefaultIPInfo todo
func NewDefaultIPInfo() *IPInfo {
	return &IPInfo{}
}

// 中国|0|四川省|成都市|电信
func ParseIpInfoFromString(raw string) (*IPInfo, error) {
	list := strings.Split(raw, "|")
	if len(list) != 5 {
		return nil, fmt.Errorf("format error")
	}
	return &IPInfo{
		Country:  list[0],
		Region:   list[1],
		Province: list[2],
		City:     list[3],
		ISP:      list[4],
	}, nil
}

// IPInfo todo
type IPInfo struct {
	Country  string `bson:"country" json:"country"`
	Region   string `bson:"region" json:"region"`
	Province string `bson:"province" json:"province"`
	City     string `bson:"city" json:"city"`
	ISP      string `bson:"isp" json:"isp"`
}

func (ip IPInfo) String() string {
	return ip.Country + "|" + ip.Region + "|" + ip.Province + "|" + ip.City + "|" + ip.ISP
}

func (ip *IPInfo) IsPublic() bool {
	return ip.City != "内网IP"
}
