package admin

import (
	"fmt"
	"github.com/dwnGnL/pg-contests/internal/application"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (ah *adminHandler) createContest(c *gin.Context) {
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

	bearerToken := c.Request.Header.Get("Authorization")
	tokenDetails, err := ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	request.CreatedBy = strconv.FormatInt(tokenDetails.ID, 10)
	contest, err := app.CreateContest(request)
	if err != nil {
		goerrors.Log().WithError(err).Error("create contest error")
		errorModel.Error.Message = "create contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contest)
}

func (ah *adminHandler) getAllContest(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	pagination := repository.GetPaginateSettings(c.Request)

	contests, err := app.GetAllContest(pagination)
	if err != nil {
		goerrors.Log().WithError(err).Error("get all contest error")
		errorModel.Error.Message = "get all contest error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	c.JSON(http.StatusOK, contests)
}

func (ah *adminHandler) getContestById(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
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

func (ah *adminHandler) deleteContestById(c *gin.Context) {
	errorModel := repository.ErrorResponse{}
	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
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

func (ah *adminHandler) updateContest(c *gin.Context) {
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

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
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

func (ah *adminHandler) changeStatus(c *gin.Context) {

	var (
		errorModel = repository.ErrorResponse{}
		msg        string
		newStatus  bool
		contestID  int64
		err        error
		app        application.Core
	)
	app, err = application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
		return
	}

	contestID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		goerrors.Log().WithError(err).Error("Parse contest id error")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	newStatus, err = app.ChangeStatus(contestID)
	if err != nil {
		goerrors.Log().WithError(err).Error("contest ChangeStatus error")
		errorModel.Error.Message = "contest ChangeStatus error: " + err.Error()
		c.JSON(http.StatusInternalServerError, errorModel)
		return
	}
	switch newStatus {
	case true:
		msg = fmt.Sprintf("Contest #%d was activated", contestID)
	case false:
		msg = fmt.Sprintf("Contest #%d was disabled", contestID)
	}

	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (ah *adminHandler) migrate(c *gin.Context) {
	errorModel := repository.ErrorResponse{}

	app, err := application.GetAppFromRequest(c)
	if err != nil {
		goerrors.Log().Warn("fatal err: %w", err)
		c.AbortWithStatus(http.StatusBadGateway)
		return
	}

	bearerToken := c.Request.Header.Get("Authorization")
	_, err = ah.jwtClient.ExtractTokenMetadata(bearerToken)
	if err != nil {
		goerrors.Log().WithError(err).Error("ExtractTokenMetadata error")
		errorModel.Error.Message = err.Error()
		c.JSON(http.StatusUnauthorized, errorModel)
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
