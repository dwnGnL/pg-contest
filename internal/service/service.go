package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/internal/repository"
)

type repositoryIter interface {
	GetAllContest() (*[]repository.Contest, error)
	GetAllContestByUserID(userID int64) (*[]repository.UserContestResp, error)
	CreateContest(contest repository.Contest) (*repository.Contest, error)
	UpdateContest(contest repository.Contest) (*repository.Contest, error)
	ChangeContestInfo(contest repository.Contest) (*repository.Contest, error)
	DeleteContest(contest repository.Contest) error
	GetContest(contestID int64) (*repository.Contest, error)
	GetUserTikets(userID, tiketID int64) (*repository.UserTickets, error)
	Migrate() error
	SubscribeContest(contest repository.Contest, userID int64) error
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

func (s ServiceImpl) CheckAndReturnContestByUserID(contestID, userID int64) (*repository.Contest, error) {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return nil, err
	}
	if !contest.Active {
		return nil, fmt.Errorf("contest not active")
	}
	tiket, err := s.repo.GetUserTikets(userID, contestID)
	if err != nil {
		return nil, err
	}

	if tiket == nil || tiket.Canseled {
		return nil, fmt.Errorf("yout tiket is canseled")
	}
	return contest, nil
}

func (s ServiceImpl) GetAllContest() (*[]repository.Contest, error) {
	contests, err := s.repo.GetAllContest()
	if err != nil {
		return nil, err
	}
	return contests, nil
}

func (s ServiceImpl) GetAllContestByUserID(userID int64) (*[]repository.UserContestResp, error) {
	contests, err := s.repo.GetAllContestByUserID(userID)
	if err != nil {
		return nil, err
	}
	return contests, nil
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

func (s ServiceImpl) ChangeStatus(contestID int64) error {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return err
	}
	_, err = s.repo.ChangeContestInfo(repository.Contest{ID: contestID, Active: !contest.Active})
	if err != nil {
		return err
	}
	return nil
}

func (s ServiceImpl) DeleteContest(contestID int64) error {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return err
	}
	err = s.repo.DeleteContest(*contest)
	if err != nil {
		return err
	}
	return nil
}

func (s ServiceImpl) Migrate() error {
	err := s.repo.Migrate()
	if err != nil {
		return err
	}
	return nil
}

func (s ServiceImpl) SubscribeContest(contestID int64, jwtToken string, userID int64) error {
	var (
		bearer string
		header map[string]string
		res    interface{}
		err    error
	)

	bearer = "Bearer " + jwtToken

	header = map[string]string{"Authorization": bearer}

	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		return err
	}
	req := struct {
		Amount float64 `json:"amount"`
	}{contest.Price}

	body, err := json.Marshal(&req)
	if err != nil {
		return err
	}
	if err = sendRequest("POST", s.conf.ApiURL, bytes.NewBuffer(body), &res, &header); err != nil {
		return err
	}
	err := s.repo.SubscribeContest(contest, userID)
	if err != nil {
		return err
	}
	return nil
}
