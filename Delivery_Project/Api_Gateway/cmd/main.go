package main

import (
	"api-gateway-service/api"
	"api-gateway-service/api/handler"
	"api-gateway-service/config"
	"api-gateway-service/pkg/logger"
	"api-gateway-service/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	fmt.Printf("config: %+v/n", cfg)

	// Setup Logger
	loggerLevel := logger.LevelDebug
	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
	case config.TestMode:
		loggerLevel = logger.LevelDebug
	default:
		loggerLevel = logger.LevelInfo
	}

	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer logger.Cleanup(log)

	grpcSrvc, err := services.NewGrpcClients(cfg)
	if err != nil {
		panic(err)
	}

	r := gin.New()

	h := handler.NewHandler(cfg, log, grpcSrvc)

	api.SetUpApi(r, h, cfg)

	fmt.Println("Start api gateway...")

	err = r.Run(cfg.HTTPPort)
	if err != nil {
		return
	}
}
