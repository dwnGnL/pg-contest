package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
)

type repositoryIter interface {
	GetAllContest(pagination *repository.Pagination) (*repository.Pagination, error)
	GetAllContestByUserID(userID int64, pagination *repository.Pagination) (*repository.Pagination, error)
	GetContestStatsById(contestID, currentQuestionID int64, pagination *repository.Pagination) (*repository.Pagination, error)
	GetContestStatsForUser(contestID, userID, currentQuestionID int64) (*repository.ContestStats, error)
	GetContestFullStatsForUser(contestID, userID int64, currentQuestionOrder int) (*repository.Contest, error)
	CreateContest(contest repository.Contest) (*repository.Contest, error)
	UpdateContest(contest repository.Contest) (*repository.Contest, error)
	ChangeContestInfo(contest *repository.Contest) error
	DeleteContest(contest repository.Contest) error
	GetContest(contestID int64) (*repository.Contest, error)
	GetContestInfo(contestID int64) (*repository.Contest, error)
	GetUserTikets(userID, tiketID int64) (*repository.UserTickets, error)
	Migrate() error
	SubscribeContest(contest repository.Contest, userID int64) error
	ContestAvailability(contestID int64, userID int64) (*repository.Contest, error)
	GetUserContest(contestID int64, userID int64) (*repository.UserContests, error)
	SubmitAnswer(userAnswer *repository.UserAnswers) (err error)
}

type ServiceImpl struct {
	conf *config.Config
	repo repositoryIter
}

type Option func(*ServiceImpl)

func New(conf *config.Config, repo repositoryIter, opts ...Option) *ServiceImpl {
	s := ServiceImpl{
		conf: conf,
		repo: repo,
	}

	for _, opt := range opts {
		opt(&s)
	}

	return &s
}

var SubscribeErr = fmt.Errorf("please subscribe contest to continue")

func (s ServiceImpl) CheckAndReturnContestByUserID(contestID, userID int64) (*repository.Contest, error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return nil, fmt.Errorf("GetContest err: %w", err)
	}
	if !*contest.Active {
		return nil, fmt.Errorf("contest not active")
	}
	userContest, err := s.repo.GetUserContest(contestID, userID)
	if err != nil {
		return nil, fmt.Errorf("GetUserContest err: %w", err)
	}

	if userContest == nil || userContest != nil && userContest.ContestID != contestID {
		return nil, SubscribeErr
	}

	return contest, nil
}

func (s ServiceImpl) CalculateTimeForQuestion(contestID, questionID int64) (resTime int64, err error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		goerrors.Log().Warnln("err on GetContest ", err)
		return
	}

	layout := "2006-01-02T15:04Z07:00"
	startTime, err := time.Parse(layout, contest.StartTime)
	if err != nil {
		goerrors.Log().Warnln("err on time.Parse ", err)
		return
	}
	startTimeUnix := startTime.Unix()

	now := time.Now().Unix() - startTimeUnix

	sort.Slice(contest.Questions, func(i, j int) bool {
		return contest.Questions[i].Order < contest.Questions[j].Order
	})

	var totalTime int64
	for _, v := range contest.Questions {
		if v.ID == questionID {
			resTime = now - totalTime
			return
		}
		totalTime += v.Time
	}
	return
}

func (s ServiceImpl) GetCurrentQuestion(contestID int64) (question repository.Question, err error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		goerrors.Log().Warnln("err on GetContest ", err)
		return
	}

	layout := "2006-01-02T15:04Z07:00"
	startTime, err := time.Parse(layout, contest.StartTime)
	if err != nil {
		goerrors.Log().Warnln("err on time.Parse ", err)
		return
	}
	startTimeUnix := startTime.Unix()

	now := time.Now().Unix() - startTimeUnix

	sort.Slice(contest.Questions, func(i, j int) bool {
		return contest.Questions[i].Order < contest.Questions[j].Order
	})

	var totalTime int64
	for _, v := range contest.Questions {
		totalTime += v.Time
		if now <= totalTime {
			question = v
			return
		}
	}
	return
}

func (s ServiceImpl) GetAllContest(pagination *repository.Pagination) (*repository.Pagination, error) {
	contests, err := s.repo.GetAllContest(pagination)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s ServiceImpl) GetAllContestByUserID(userID int64, pagination *repository.Pagination) (*repository.Pagination, error) {
	contests, err := s.repo.GetAllContestByUserID(userID, pagination)
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s ServiceImpl) GetContestStatsById(contestID int64, pagination *repository.Pagination) (*repository.Pagination, error) {
	currentQuestion, err := s.GetCurrentQuestion(contestID)
	if err != nil {
		return nil, err
	}
	contestStats, err := s.repo.GetContestStatsById(contestID, currentQuestion.ID, pagination)
	if err != nil {
		return nil, err
	}
	return contestStats, nil
}

func (s ServiceImpl) GetContestStatsForUser(contestID, userID int64) (*repository.ContestStats, error) {
	currentQuestion, err := s.GetCurrentQuestion(contestID)
	if err != nil {
		return nil, err
	}
	contestStats, err := s.repo.GetContestStatsForUser(contestID, userID, currentQuestion.ID)
	if err != nil {
		return nil, err
	}
	return contestStats, nil
}

func (s ServiceImpl) GetContestFullStatsForUser(contestID, userID int64) (*repository.Contest, error) {
	currentQuestion, err := s.GetCurrentQuestion(contestID)
	if err != nil {
		return nil, err
	}
	contest, err := s.repo.GetContestFullStatsForUser(contestID, userID, currentQuestion.Order)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s ServiceImpl) GetContest(contestID int64) (*repository.Contest, error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (s ServiceImpl) CreateContest(contest repository.Contest) (*repository.Contest, error) {
	createdContest, err := s.repo.CreateContest(contest)
	if err != nil {
		return nil, err
	}
	return createdContest, nil
}

func (s ServiceImpl) UpdateContest(contest repository.Contest) (*repository.Contest, error) {
	updatedContest, err := s.repo.UpdateContest(contest)
	if err != nil {
		return nil, err
	}
	return updatedContest, nil
}

func (s ServiceImpl) ChangeStatus(contestID int64) (newStatus bool, err error) {
	contest, err := s.repo.GetContestInfo(contestID)
	if err != nil {
		return
	}
	*contest.Active = !*contest.Active
	return *contest.Active, s.repo.ChangeContestInfo(contest)
}

func (s ServiceImpl) DeleteContest(contestID int64) (err error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return
	}
	return s.repo.DeleteContest(*contest)
}

func (s ServiceImpl) Migrate() error {
	return s.repo.Migrate()
}

func (s ServiceImpl) SubmitAnswer(userAnswer *repository.UserAnswers) (err error) {
	return s.repo.SubmitAnswer(userAnswer)
}

func (s ServiceImpl) SubscribeContest(userContest *repository.UserContests, jwtToken string) error {
	var (
		// header map[string]string
		// res    interface{}
		err error
	)

	// header = map[string]string{"Authorization": jwtToken}

	contest, err := s.repo.ContestAvailability(userContest.ContestID, userContest.UserID)
	if err != nil {
		return err
	}
	// req := struct {
	// 	Amount float64 `json:"amount"`
	// }{contest.Price}

	// body, err := json.Marshal(&req)
	// if err != nil {
	// 	return err
	// }

	goerrors.Log().Info("contest:", contest)
	// if err = s.SendRequest("POST", bytes.NewBuffer(body), &res, &header); err != nil {
	// 	return err
	// }
	err = s.repo.SubscribeContest(*contest, userContest.UserID)
	if err != nil {
		return err
	}
	return nil
}
