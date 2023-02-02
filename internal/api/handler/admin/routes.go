package admin

import (
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/lib/token"
	"github.com/gin-gonic/gin"
)

type adminHandler struct {
	jwtClient token.JwtToken[AdminAccessDetails]
}

func newAdminHandler(cfg *config.Config) *adminHandler {
	return &adminHandler{
		jwtClient: token.New[AdminAccessDetails](cfg.AdminPrivKey),
	}
}

func GenRouting(r *gin.RouterGroup, cfg *config.Config) {
	admin := newAdminHandler(cfg)

	//admin
	r.POST("/contest", admin.createContest)
	r.GET("/contests", admin.getAllContest)
	r.GET("/contest/:id", admin.getContestById)
	r.GET("/contest/:id/changeStatus", admin.changeStatus)
	r.DELETE("/contest/:id", admin.deleteContestById)
	r.PUT("/contest", admin.updateContest)
	r.POST("/migrate", admin.migrate)

}
