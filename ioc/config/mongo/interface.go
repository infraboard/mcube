package mongo

import (
	"github.com/infraboard/mcube/ioc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGODB = "mongodb"
)

type ClientGetter interface {
	// 获取client, 主要用于开始session
	Client() *mongo.Client
	// 获取DB
	GetDB() *mongo.Database
}

func GetClientGetter() ClientGetter {
	return ioc.Config().Get(MONGODB).(ClientGetter)
}
