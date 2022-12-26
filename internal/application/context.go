package application

import (
	"context"
	"errors"

	"github.com/gin-gonic/gin"
)

type ContextValue int

const (
	ContextApp ContextValue = iota
)

func GetAppFromContext(ctx context.Context) (Core, error) {
	if app, ok := ctx.Value(ContextApp).(Core); ok {
		return app, nil
	}

	return nil, errors.New("cannot get app container from request")
}

func GetAppFromRequest(r *gin.Context) (Core, error) {
	return GetAppFromContext(r.Request.Context())
}
