package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dwnGnL/pg-contests/internal/api/wshandler"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/lib/goerrors"

	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/gin-gonic/gin"
)

type GracefulStopFuncWithCtx func(ctx context.Context) error

func SetupHandlers(core application.Core, cfg *config.Config) GracefulStopFuncWithCtx {
	c := gin.New()
	c.Use(application.WithApp(core), application.WithCORS())
	apiv1 := c.Group("/api/v1/")
	// apiv1.Use() добавить проверку токена
	generateAPIV1Routing(apiv1)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.ListenPort),
		Handler: c,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			goerrors.Log().Fatalf("listen: %s\n", err)
		}
	}()
	return srv.Shutdown
}

func generateAPIV1Routing(gE *gin.RouterGroup) {
	wshandler.GenRouting(gE)
}
