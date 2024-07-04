package main

import (
	"context"
	"fmt"
	"github.com/dbacilio88/go-application-api/config"
	"github.com/dbacilio88/go-application-api/config/log"
	"github.com/dbacilio88/go-application-api/pkg/controllers"
	"github.com/dbacilio88/go-application-api/pkg/models/response"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

func main() {

	config.LoadConfiguration("./config/")
	logger, _ := log.LoggerConfiguration(config.Configuration.Log.Level)
	defer closeLoggerHandler()(logger)

	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	logger.Info("starting server main...")

	var srvCfg controllers.Config
	server := controllers.NewServer(logger, &srvCfg)
	stopCh := SetupSignalHandler()

	RequestMiddleware()
	go server.ListenAndServe(stopCh)

	response.AppHealthState.Status = response.HealthStatusStarting
	response.AppHealthState.Readiness = response.HealthValueStatus{Status: response.HealthStatusUp}
	response.AppHealthState.Liveliness = response.HealthValueStatus{Status: response.HealthStatusUp}

	server.SetConnectionHealth()
	select {}
}

func RequestMiddleware() {
	log.Reset()
	requestId := uuid.New()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", requestId)
	log.LoggerInstance = log.LoggerInstance.With(zap.String("requestId", requestId.String()))
	log.WithCtx(ctx, log.LoggerInstance)
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
