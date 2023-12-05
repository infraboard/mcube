package mongo

import (
	"github.com/infraboard/mcube/v2/ioc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	AppName = "mongodb"
)

func DB() *mongo.Database {
	return ioc.Config().Get(AppName).(*mongoDB).GetDB()
}

func Client() *mongo.Client {
	return ioc.Config().Get(AppName).(*mongoDB).Client()
}
