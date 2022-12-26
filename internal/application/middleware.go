package application

import (
	"context"

	"github.com/gin-gonic/gin"
)

func WithApp(app Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ContextApp, app)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
