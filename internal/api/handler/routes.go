package handler

import (
	"github.com/gin-gonic/gin"
)

func GenRouting(r *gin.RouterGroup) {
	r.POST("/contest", createContest)
	r.GET("/contests", getAllContest)
	r.GET("/user/contests", getAllContestByUserID)
	r.GET("/contest/:id", getContestById)
	r.GET("/contest/:id/changeStatus", changeStatus)
	r.DELETE("/contest/:id", deleteContestById)
	r.POST("/contest/:id/subscribe", subscribeContestById)
	r.PUT("/contest", updateContest)
	r.POST("/migrate", migrate)
}
