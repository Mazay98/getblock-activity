package main

import (
	"context"
	"errors"
	"log"
	"os"

	"getBlock/internal/config"
	"getBlock/internal/environment"
	"getBlock/internal/service/blockio"
	ll "getBlock/pkg/logger"
	"github.com/chapsuk/grace"
	"go.uber.org/zap"
)

//nolint:gochecknoglobals
var (
	version   = "unknown"
	buildTime = "unknown"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		if errors.Is(err, config.ErrHelp) {
			os.Exit(0)
		}
		log.Fatalf("failed to read app config: %v", err)
	}

	logger, err := ll.New(version, cfg.Env, cfg.Logger.Level)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
	defer logger.Sync() //nolint:errcheck

	ctx := grace.ShutdownContext(context.Background())
	ctx = environment.CtxWithEnv(ctx, cfg.Env)
	ctx = environment.CtxWithVersion(ctx, version)
	ctx = environment.CtxWithBuildTime(ctx, buildTime)

	bService, err := blockio.New(ctx, logger, &cfg.Blockio)
	if err != nil {
		logger.Fatal("Failed create a blockio service", zap.Error(err))
		os.Exit(1)
	}

	bService.GetTopActivity(ctx, cfg.Positions)
}
