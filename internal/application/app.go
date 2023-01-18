package application

import (
	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/repository"
)

type Core interface {
	GetAllContest() (*[]repository.Contest, error)
	GetContest(contestID int64) (*repository.Contest, error)
	DeleteContest(contestID int64) error
	CreateContest(contest repository.Contest) (*repository.Contest, error)
	UpdateContest(contest repository.Contest) (*repository.Contest, error)
	ChangeStatus(contestID int64) error
	CheckAndReturnContestByUserID(contestID, userID int64) (*repository.Contest, error)
	GenerateAndProcessChan(contestID int64) <-chan models.WsResponse
	Generate(contestID int64) models.WsResponse
	Migrate() error
}
