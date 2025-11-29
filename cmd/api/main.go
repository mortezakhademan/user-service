package main

import (
	"context"
	"flag"
	"git.ramooz.org/ramooz/golang-components/logger"
	"git.ramooz.org/ramooz/golang-components/microservice/connection"
	"github.com/joho/godotenv"
	"github.com/mortezakhademan/user-service-sample/internal/config"
	"github.com/mortezakhademan/user-service-sample/internal/db"
	"github.com/mortezakhademan/user-service-sample/internal/transport"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	setTimeZone()
	initLogger()
	setEnvFile()
	if err := initDatabase(ctx); err != nil {
		config.Logger.Fatal(err)
	}
	grpcServer, err := transport.InitGrpcService(config.GetHttpAddress(), config.GetGrpcPort())
	if err != nil {
		config.Logger.Fatal(err)
	}
	httpServer, err := transport.InitRestService(ctx, config.GetHttpAddress(), config.GetHttpPort())
	if err != nil {
		config.Logger.Fatal(err)
	}

	config.Logger.Fatal("Grace fully shutdown server: ", connection.GracefullyShutdownServer(ctx, grpcServer, httpServer))

}

func setEnvFile() {
	filePath := flag.String("env", "", "env file path")
	flag.Parse()
	if filePath == nil || *filePath == "" {
		config.Logger.Fatal("env file path is required, pass as argument. example: go run main.go -env /path1/.env")
	}
	if err := godotenv.Load(*filePath); err != nil {
		config.Logger.Fatal("failed to load env file: ", err)
	}
}

func initLogger() error {
	log, err := logger.NewLogger(1, "user",
		&logger.Options{
			LogLevel:      logger.DebugLevel,
			ConsoleWriter: true,
			TimeFormat:    logger.RFC3339_TIME,
			Colorable:     true,
			Development:   true,
			LogPath:       "",
		})
	if err != nil {
		return err
	}
	config.Logger = log
	return nil
}

func initDatabase(ctx context.Context) error {
	mongoClient, err := db.ConnectMongo(ctx, config.GetMongoUri())
	if err != nil {
		log.Fatal(err)
	}

	database := mongoClient.Database("sample_project")

	config.MongoClient = mongoClient
	config.DB = database
	return nil
}

func setTimeZone() {
	os.Setenv("TZ", "Asia/Tehran")
}
