package mongo

import (
	"github.com/infraboard/mcube/v2/ioc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AppName = "mongo"
)

func DB() *mongo.Database {
	return Get().GetDB()
}

func Client() *mongo.Client {
	return Get().Client()
}

func Get() *mongoDB {
	obj := ioc.Config().Get(AppName)
	if obj == nil {
		return defaultConfig
	}
	return obj.(*mongoDB)
}
