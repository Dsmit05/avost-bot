package main

import (
	"context"
	"github.com/Dsmit05/avost-bot/pkg/stat"
	"os/signal"
	"syscall"

	"github.com/Dsmit05/avost-bot/internal/api"
	"github.com/Dsmit05/avost-bot/internal/api/cryptography"
	"github.com/Dsmit05/avost-bot/internal/bot"
	"github.com/Dsmit05/avost-bot/internal/config"
	"github.com/Dsmit05/avost-bot/internal/logger"
	"github.com/Dsmit05/avost-bot/internal/repository/cache"
	"github.com/Dsmit05/avost-bot/pkg/animevost/nclient"
	"github.com/pkg/errors"

	"os"
	"runtime/debug"
	"time"
)

// @title        avost-bot
// @version      1.0.0
// @description  API bot

// @host      localhost:8080
// @basePath  /api/v1/

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorizations
func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("panic", errors.New(string(debug.Stack())))
		}
	}()

	// Init logger
	if err := logger.InitLogger(false, "data/logs.json"); err != nil {
		panic(err)
	}

	// Init Config
	cfg := config.NewAnimeBotConfg(5*time.Second, 15*time.Minute, ":8080")

	logger.Info("BuildVersion", cfg.BuildVersion)

	// Init animevost client
	cli, err := nclient.NewClient("https://api.animevost.org/v1")
	if err != nil {
		logger.Error("animevost.NewClient()", err)
		return
	}

	logger.Info("animevost.NewClient()", "start")

	// Cache for users
	fcache := cache.NewFavouritesCache("data/db", 2*time.Minute)
	defer func() {
		if err = fcache.Stop(); err != nil {
			logger.Error("fcache.Stop()", err)
		}
	}()
	logger.Info("cache.NewFavouritesCache()", "created")

	go func() {
		if err = fcache.Start(context.Background()); err != nil {
			logger.Error("fcache.Start()", err)
		}
	}()
	logger.Info("fcache.Start()", "started")

	// Statistic
	st, err := stat.NewStatistic("./data/stat")
	if err != nil {
		logger.Error("stat.NewStatistic(./data/stat)", err)
		return
	}
	defer st.Stop()

	// Start main logic bot
	aBot, err := bot.NewBot(cfg, cli, fcache, st)
	if err != nil {
		logger.Error("bot.NewBot()", err)
		return
	}

	go aBot.InitControllersV1().Start()
	defer aBot.Stop()

	logger.Info("aBot.InitControllersV1().Start()", "started")

	// Start api
	tokenJWT := cryptography.NewTokenJWT(cfg.APIConfig.SecretKeyJWT)

	apiServer := api.NewServer(fcache, aBot, tokenJWT, st, cfg.APIConfig.Addr).AddV1("/api/v1").AddSwagger(cfg.Address)
	go apiServer.Start()

	// Graceful Stop
	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("stop app", "stop with os signal")
}
