package application

import (
	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/repository"
)

type Core interface {
	GetAllContest(pagination *repository.Pagination) (*repository.Pagination, error)
	GetAllContestByUserID(userID int64, pagination *repository.Pagination) (*repository.Pagination, error)
	GetContestStatsById(contestID int64, pagination *repository.Pagination) (*repository.Pagination, error)
	GetContestStatsForUser(contestID, userID int64) (*repository.ContestStats, error)
	GetContestFullStatsForUser(contestID, userID int64) (*repository.Contest, error)
	GetContest(contestID int64) (*repository.Contest, error)
	DeleteContest(contestID int64) error
	CreateContest(contest repository.Contest) (*repository.Contest, error)
	UpdateContest(contest repository.Contest) (*repository.Contest, error)
	ChangeStatus(contestID int64) (newStatus bool, err error)
	CheckAndReturnContestByUserID(contestID, userID int64) (*repository.Contest, error)
	GenerateAndProcessChan(contestID int64) <-chan models.WsResponse
	Generate(contestID int64) models.WsResponse
	Migrate() error
	SubscribeContest(userContest *repository.UserContests, jwtToken string) error
	CalculateTimeForQuestion(contestID, questionID int64) (int64, error)
	GetCurrentQuestion(contestID int64) (repository.Question, error)
	SubmitAnswer(userAnswer *repository.UserAnswers) (err error)
}
