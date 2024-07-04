package main

import (
	"fmt"
	"github.com/dbacilio88/go-application-api/config"
	"github.com/dbacilio88/go-application-api/config/log"
	"github.com/dbacilio88/go-application-api/pkg/models/response"
	"github.com/dbacilio88/go-application-api/pkg/node"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {
	config.LoadConfigurationMicroservice("./config/")
	logger, _ := log.ApplyLoggerConfiguration(config.Configuration.Log.Level)
	defer closeLoggerHandler()(logger)
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	logger.Info("starting node...")
	var srvCfg node.Config
	server := node.NewServer(logger, &srvCfg)
	stopCh := SetupSignalHandler()
	go server.ListenAndServe(stopCh)
	response.AppHealthState.Status = response.HealthStatusStarting
	response.AppHealthState.Readiness = response.HealthValueStatus{Status: response.HealthStatusUp}

	server.SetConnectionHealth()
	select {}
}

func closeLoggerHandler() func(logger *zap.Logger) {
	return func(logger *zap.Logger) {
		if err := logger.Sync(); err != nil {
			fmt.Println(err)
		}
	}
}

var onlyOneSignalHandler = make(chan struct{})

func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler)
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1)
	}()
	return stop
}
