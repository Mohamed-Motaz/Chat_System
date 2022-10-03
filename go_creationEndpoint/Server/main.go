package main

import (
	logger "Server/Logger"
	s "Server/Server"
	utils "Server/Utils"

	"os"
	"os/signal"
	"syscall"
)

const LOG_FOLDER = "ServerLogs"

func main() {

	utils.InitDirectories(LOG_FOLDER)
	logger.InitLogger(LOG_FOLDER)

	s, err := s.NewServer()
	if err != nil {
		logger.FailOnError(logger.SERVER, logger.ESSENTIAL, "Error while creating server %v", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-signalCh //block until user exits
	logger.LogInfo(logger.SERVER, logger.ESSENTIAL, "Received a quit sig %+v\nCleaning up resources. Goodbye", sig)
	s.Mq.Close()
}
