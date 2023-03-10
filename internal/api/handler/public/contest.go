package public

import (
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ph *publicHandler) getAllContestByUserID(c *gin.Context) {
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
	bearerToken := c.Request.Header.Get("Authorization")

	tokenDetails, err := ph.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	pagination := repository.GetPaginateSettings(c.Request)

	contests, err := app.GetAllContestByUserID(tokenDetails.ID, pagination)
	if err != nil {
		goerrors.Log().WithError(err).Error("get all contest by userID error")
		errorModel.Error.Message = "get all contest by userID error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contests)
}

func (ph *publicHandler) subscribeContestById(c *gin.Context) {
	var (
		app        application.Core
		errorModel = repository.ErrorResponse{}
		jwtToken   string
		err        error
		request    = repository.UserContests{}
	)

	if err = c.ShouldBindJSON(&request); err != nil {
		goerrors.Log().WithError(err).Error("bind request error")
		errorModel.Error.Message = "bind request error: " + err.Error()
		c.JSON(http.StatusBadRequest, errorModel)
		return
	}

	app, err = application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	jwtToken = c.Request.Header.Get("Authorization")

	bearerToken := c.Request.Header.Get("Authorization")
	tokenDetails, err := ph.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}
	request.UserID = tokenDetails.ID

	err = app.SubscribeContest(&request, jwtToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("subscribe contest error")
		errorModel.Error.Message = "subscribe contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success"})
}

func (ph *publicHandler) getContestStatsById(c *gin.Context) {
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

	pagination := repository.GetPaginateSettings(c.Request)
	pagination.Sort = ""
	contestStats, err := app.GetContestStatsById(contestID, pagination)
	if err != nil {
		goerrors.Log().WithError(err).Error("get contest stats by contestID error")
		errorModel.Error.Message = "get contest stats by contestID error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contestStats)
}

func (ph *publicHandler) getContestStatsForUser(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	tokenDetails, err := ph.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}
	userID := tokenDetails.ID

	contestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	contestStatsForUser, err := app.GetContestStatsForUser(contestID, userID)
	if err != nil {
		goerrors.Log().WithError(err).Error("get contest stats by contestID error")
		errorModel.Error.Message = "get contest stats by contestID error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contestStatsForUser)
}

func (ph *publicHandler) getContestFullStatsForUser(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	tokenDetails, err := ph.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}
	userID := tokenDetails.ID

	contestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	contestStatsForUser, err := app.GetContestFullStatsForUser(contestID, userID)
	if err != nil {
		goerrors.Log().WithError(err).Error("get contest full stats by contestID error")
		errorModel.Error.Message = "get contest full stats by contestID error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contestStatsForUser)
}
