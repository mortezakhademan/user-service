package config

import (
	"git.ramooz.org/ramooz/golang-components/logger"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"os"
)

var (
	Logger      *logger.LogService
	MongoClient *mongo.Client
	DB          *mongo.Database
)

func GetHttpAddress() string {
	return os.Getenv("HTTP_ADDRESS")
}

func GetHttpPort() string {
	return os.Getenv("HTTP_PORT")
}

func GetGrpcPort() string {
	return os.Getenv("GRPC_PORT")
}

func GetMongoUri() string {
	return os.Getenv("MONGO_URI")
}
