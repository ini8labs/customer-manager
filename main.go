package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ini8labs/lsdb"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/ini8labs/customer-manager/docs"
	"github.com/ini8labs/customer-manager/src/apis"
)

// @title My API
// @version 1.0
// @description This is Lottery Project API
// @host localhost:3000
// @BasePath /api/v1
// @schemes http
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err.Error())
	}
	_, ok := os.LookupEnv("MONGO_DB_CONN_STRING")
	if !ok {
		panic("value not found")
	}

	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	dbClient, err := lsdb.NewClient()
	if err != nil {
		panic(err.Error())
	}

	if err := dbClient.OpenConnection(); err != nil {
		panic(err.Error())
	}
	defer dbClient.CloseConnection()

	server := apis.Server{
		Logger: logger,
		Client: dbClient,
		Addr:   ":3000",
	}

	go func() {
		if err := apis.NewServer(server); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt
	logger.Info("Closing the Server")

}
