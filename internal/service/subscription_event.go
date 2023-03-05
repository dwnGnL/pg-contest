package service

import (
	"sort"
	"time"

	"github.com/dwnGnL/pg-contests/internal/api/models"
	"github.com/dwnGnL/pg-contests/internal/repository"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
)

func (s ServiceImpl) GenerateAndProcessChan(contestID int64) <-chan models.WsResponse {
	ch := make(chan models.WsResponse)
	go s.chanWorker(ch, contestID)
	return ch
}
func (s ServiceImpl) Generate(contestID int64) models.WsResponse {
	contest, err := s.repo.GetContest(contestID)
	if err != nil {
		goerrors.Log().Warnln("err on GetContest ", err)
		return models.WsResponse{}
	}
	var resp models.WsResponse
	resp.TotalStep = len(contest.Questions)

	layout := "2006-01-02T15:04Z07:00"
	startTime, err := time.Parse(layout, contest.StartTime)
	if err != nil {
		goerrors.Log().Warnln("err on time.Parse ", err)
		return models.WsResponse{}
	}
	startTimeUnix := startTime.Unix()
	now := time.Now().Unix()
	remained := now - startTimeUnix
	if remained < 0 {
		resp.TotalTime = startTimeUnix
		resp.CountDown = (remained) * -1
		resp.ContestStatus = models.Waiting
		return resp
	}
	sort.Slice(contest.Questions, func(i, j int) bool {
		return contest.Questions[i].Order < contest.Questions[j].Order
	})
	var totalTime int64
	for i, v := range contest.Questions {
		totalTime += v.Time
		if now-totalTime >= startTimeUnix {
			resp.Questions = append(resp.Questions, convertRepQToWsQ(contest.Questions[i]))
			continue
		}
		break
	}

	countPassed := len(resp.Questions)
	if countPassed != len(contest.Questions) {
		resp.Questions = append(resp.Questions, convertRepQToWsQ(contest.Questions[countPassed]))
		resp.ActiveQuestionID = contest.Questions[countPassed].ID
		resp.Step = countPassed + 1
		resp.TotalTime = contest.Questions[countPassed].Time
	}

	resp.ContestStatus = models.Start
	resp.CountDown = (totalTime + startTimeUnix) - now
	resp.TotalStep = len(contest.Questions)
	if countPassed == resp.TotalStep && resp.CountDown <= 0 {
		resp.ContestStatus = models.End
	}
	if resp.CountDown == 0 {
		resp.CountDown++
	}
	if resp.CountDown < 0 {
		resp.ContestStatus = models.End
	}
	return resp
}

func (s ServiceImpl) chanWorker(ch chan<- models.WsResponse, contestID int64) {
	for {
		resp := s.Generate(contestID)
		truePointer := true
		if resp.ContestStatus == models.End {
			s.repo.ChangeContestInfo(&repository.Contest{ID: contestID, IsEnd: &truePointer})
			ch <- resp
			close(ch)
			return
		}
		ch <- resp
		time.Sleep(time.Duration(resp.CountDown) * time.Second)

	}

}

func convertRepQToWsQ(question repository.Question) models.WsQuestion {
	return models.WsQuestion{
		ID:      question.ID,
		Order:   question.Order,
		Title:   question.Title,
		Answers: convertRepAToWsA(question.Answers),
	}
}
func convertRepAToWsA(answer []repository.Answer) []models.WsAnswer {
	var wsAnswer []models.WsAnswer
	for _, v := range answer {
		wsAnswer = append(wsAnswer, models.WsAnswer{
			ID:    v.ID,
			Title: v.Title,
		})
	}
	return wsAnswer
}
