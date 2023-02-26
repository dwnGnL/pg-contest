package public

import (
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/lib/cachemap"
	"github.com/dwnGnL/pg-contests/lib/token"
	"github.com/gin-gonic/gin"
)

type publicHandler struct {
	contestMap *cachemap.CacheMaper[int64, *subscribeSwitcher]
	jwtClient  token.JwtToken[PublicAccessDetails]
}

func newPublicHandler(cfg *config.Config) *publicHandler {
	return &publicHandler{
		contestMap: cachemap.NewCacheMap[int64, *subscribeSwitcher](),
		jwtClient:  token.New[PublicAccessDetails](cfg.PublicPrivKey),
	}
}

func GenRouting(r *gin.RouterGroup, cfg *config.Config) {

	public := newPublicHandler(cfg)

	//user
	r.GET("/user/contests", public.getAllContestByUserID)
	r.POST("/contest/:id/subscribe", public.subscribeContestById)

	//ws
	r.Any("/connect/:contestID", public.wsContest)
}
