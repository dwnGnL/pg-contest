package handler

import (
	"fmt"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/internal/cmd"
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func createContest(c *gin.Context) {
	var request repository.Contest
	if err := c.ShouldBindJSON(&request); err != nil {
		goerrors.Log().WithError(err).Error("bind request error")
		c.JSON(http.StatusBadRequest, gin.H{"message": "bind request error: " + err.Error()})
		return
	}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	//todo: request.CreatedBy = user_id
	contest, err := app.CreateContest(request)
	if err != nil {
		goerrors.Log().WithError(err).Error("create contest error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "create contest error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, contest)
}

func getAllContest(c *gin.Context) {

	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	//todo: request.CreatedBy = user_id
	contests, err := app.GetAllContest()
	if err != nil {
		goerrors.Log().WithError(err).Error("get all contest error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "get all contest error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, contests)
}

func getContestById(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	contestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	contest, err := app.GetContest(contestID)
	if err != nil {
		goerrors.Log().WithError(err).Error("get contest error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "get contest error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, contest)
}

func deleteContestById(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	contestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = app.DeleteContest(contestID)
	if err != nil {
		goerrors.Log().WithError(err).Error("delete contest error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "delete contest error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func updateContest(c *gin.Context) {
	var request repository.Contest
	if err := c.ShouldBindJSON(&request); err != nil {
		goerrors.Log().WithError(err).Error("bind request error")
		c.JSON(http.StatusBadRequest, gin.H{"message": "bind request error: " + err.Error()})
		return
	}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	contest, err := app.UpdateContest(request)
	if err != nil {
		goerrors.Log().WithError(err).Error("update contest error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "update contest error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, contest)
}

func changeStatus(c *gin.Context) {
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	contestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = app.ChangeStatus(contestID)
	if err != nil {
		goerrors.Log().WithError(err).Error("contest ChangeStatus error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "contest ChangeStatus error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func migrate(c *gin.Context) {

	cfg := config.FromFile("./config.yaml")
	fmt.Println(cfg)
	err := cmd.StartMigrate(cfg)
	if err != nil {
		goerrors.Log().WithError(err).Error("migrate error")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "migrate error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
