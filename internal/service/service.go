package service

import (
	"fmt"

	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/internal/repository"
)

type repositoryIter interface {
	GetContest(contestID int64) (*repository.Contest, error)
	GetUserTikets(userID, tiketID int64) (*repository.UserTickets, error)
	UpdateContest(contestID int64, contest repository.Contest) error
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
