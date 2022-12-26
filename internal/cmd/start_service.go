package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dwnGnL/pg-contests/internal/api"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/internal/service"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"golang.org/x/sync/errgroup"
)

const (
	gracefulStop = 5 * time.Second
)

func StartService(ctx context.Context, cfg *config.Config) error {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	s, err := buildService(ctx, cfg)
	if err != nil {
		return fmt.Errorf("build service err:%w", err)
	}
	httpgrpcGracefulStopWithCtx := api.SetupHandlers(s, cfg)
	var group errgroup.Group

	group.Go(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		goerrors.Log().Debug("wait for Ctrl-C")
		<-sigCh
		goerrors.Log().Debug("Ctrl-C signal")
		cancelCtx()
		shutdownCtx, shutdownCtxFunc := context.WithDeadline(ctx, time.Now().Add(gracefulStop))
		defer shutdownCtxFunc()

		_ = httpgrpcGracefulStopWithCtx(shutdownCtx)
		return nil
	})

	if err := group.Wait(); err != nil {
		goerrors.Log().WithError(err).Error("Stopping service with error")
	}
	return nil
}

func StartMigrate(conf *config.Config) error {
	repo, err := repository.NewRepository(conf)
	if err != nil {
		return fmt.Errorf("new repository err:%w", err)
	}
	return repo.Migrate()
}

func buildService(ctx context.Context, conf *config.Config) (application.Core, error) {
	repo, err := repository.NewRepository(conf)
	if err != nil {
		return nil, fmt.Errorf("new repository err:%w", err)
	}
	return service.New(conf, repo), nil
}
