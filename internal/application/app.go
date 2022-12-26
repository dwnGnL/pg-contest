package application

import (
	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/repository"
)

type Core interface {
	CheckAndReturnContestByUserID(contestID, userID int64) (*repository.Contest, error)
	GenerateAndProcessChan(contestID int64) <-chan models.WsResponse
	Generate(contestID int64) models.WsResponse
}
