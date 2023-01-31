package handler

import (
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/internal/middleware"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func createContest(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	var request repository.Contest
	if err := c.ShouldBindJSON(&request); err != nil {
		goerrors.Log().WithError(err).Error("bind request error")
		errorModel.Error.Message = "bind request error: " + err.Error()
		c.JSON(http.StatusBadRequest, errorModel)
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
		errorModel.Error.Message = "create contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contest)
}

func getAllContest(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
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
		errorModel.Error.Message = "get all contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contests)
}

func getAllContestByUserID(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	/*
		userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}*/

	tokenDetails, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	contests, err := app.GetAllContestByUserID(tokenDetails.ID)
	if err != nil {
		goerrors.Log().WithError(err).Error("get all contest by userID error")
		errorModel.Error.Message = "get all contest by userID error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contests)
}

func getContestById(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
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
		errorModel.Error.Message = "get contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contest)
}

func deleteContestById(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
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
		errorModel.Error.Message = "delete contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func subscribeContestById(c *gin.Context) {
	var (
		app        application.Core
		errorModel = repository.ErrorResponse{}
		jwtToken   string
		contestID  int64
		err        error
	)
	app, err = application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	jwtToken = c.Request.Header.Get("Authorization")

	tokenDetails, err := middleware.ExtractTokenMetadata(c)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	contestID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	err = app.SubscribeContest(contestID, jwtToken, tokenDetails.ID)
	if err != nil {
		goerrors.Log().WithError(err).Error("subscribe contest error")
		errorModel.Error.Message = "subscribe contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func updateContest(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	var request repository.Contest
	if err := c.ShouldBindJSON(&request); err != nil {
		goerrors.Log().WithError(err).Error("bind request error")
		errorModel.Error.Message = "bind request error: " + err.Error()
		c.JSON(http.StatusBadRequest, errorModel)
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
	errorModel := repository.ErrorResponse{}
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
		errorModel.Error.Message = "contest ChangeStatus error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func migrate(c *gin.Context) {
	errorModel := repository.ErrorResponse{}

	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}
	err = app.Migrate()
	if err != nil {
		goerrors.Log().WithError(err).Error("migrate error")
		errorModel.Error.Message = "migrate error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}
