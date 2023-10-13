package mongo

import (
	"github.com/infraboard/mcube/ioc"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGODB = "mongodb"
)

func DB() *mongo.Database {
	return ioc.Config().Get(MONGODB).(*mongoDB).GetDB()
}

func Client() *mongo.Client {
	return ioc.Config().Get(MONGODB).(*mongoDB).Client()
}
