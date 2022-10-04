package main

import (
	logger "Worker/Logger"
	utils "Worker/Utils"
	s "Worker/Worker"

	"os"
	"os/signal"
	"syscall"
)

const LOG_FOLDER = "WorkerLogs"

func main() {

	utils.InitDirectories(LOG_FOLDER)
	logger.InitLogger(LOG_FOLDER)

	s, err := s.NewWorker()
	if err != nil {
		logger.FailOnError(logger.WORKER, logger.ESSENTIAL, "Error while creating worker %v", err)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-signalCh //block until user exits
	logger.LogInfo(logger.WORKER, logger.WORKER, "Received a quit sig %+v\nCleaning up resources. Goodbye", sig)
	s.Mq.Close()
}
