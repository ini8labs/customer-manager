package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"

	"Lottery-project/src/apis"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	go func() {
		if err := apis.NewServer(":3000", logger); err != nil {
			panic(err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	logger.Info("Closing the Server")
}
