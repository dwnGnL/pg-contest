package repository

import (
	"errors"
	"fmt"
	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"gorm.io/gorm"
	"time"
)

var layout = "2006-01-02T15:04"

type Contest struct {
	ID           int64      `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	Title        string     `json:"title" binding:"required" gorm:"column:title"`
	Price        float64    `json:"price" binding:"required" gorm:"column:price"`
	PlayersCount int64      `json:"players_count" binding:"required" gorm:"players_count"`
	StartTime    string     `json:"start_time" binding:"required" gorm:"column:start_time"`
	CreatedBy    string     `json:"created_by" gorm:"column:created_by"`
	Medias       []Media    `json:"medias" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE;"`
	Questions    []Question `json:"questions" gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE"`
	Active       bool       `json:"active" gorm:"column:active;default:false"`
	IsEnd        bool       `json:"is_end" gorm:"column:is_end;default:false"`
}

type Question struct {
	ID        int64    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	ContestID int64    `json:"contest_id" binding:"required" gorm:"column:contest_id"`
	Answers   []Answer `json:"answers" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE;"`
	Title     string   `json:"title" binding:"required" gorm:"column:title"`
	Score     int      `json:"score" binding:"required" gorm:"column:score"`
	Order     int      `json:"order" binding:"required" gorm:"column:sort_order"`
	Time      int64    `json:"time" binding:"required" gorm:"column:time"`
}

type Answer struct {
	ID         int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	QuestionID int64  `json:"question_id" binding:"required" gorm:"column:contest_id"`
	Title      string `json:"title" binding:"required" gorm:"column:title"`
	IsCorrect  bool   `json:"is_correct" gorm:"column:is_correct;default:false"`
}

type Media struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Path      string `json:"path" binding:"required" gorm:"column:path"`
	OwnerID   int64  `json:"owner_id" gorm:"column:owner_id"`
	OwnerType string `json:"owner_type" gorm:"column:owner_type"`
}

type UserTickets struct {
	UserID    int64 `gorm:"column:user_id"`
	ContestID int64 `gorm:"column:contest_id"`
	Canseled  bool  `gorm:"column:canseled;default:false"`
}

func (c *Contest) Validate() error {
	_, err := time.Parse(layout, c.StartTime)
	if err != nil {
		goerrors.Log().Warnln("err on contest startTime Parse ", err)
		return err
	}
	for i, question := range c.Questions {
		if len(question.Answers) < 1 {
			err = errors.New(fmt.Sprintf("Попытка добавления вопроса №%d без ответа", i))
			goerrors.Log().Warnln(err)
			return err
		}
		correctAnsExist := false
		for _, answer := range question.Answers {
			if answer.IsCorrect {
				correctAnsExist = true
				break
			}
		}
		if !correctAnsExist {
			err = errors.New(fmt.Sprintf("Попытка добавления вопроса №%d c без правильных ответов", i))
			goerrors.Log().Warnln(err)
			return err
		}
	}
	return nil
}

func (c *Contest) Started() (bool, error) {
	startTime, err := time.Parse(layout, c.StartTime)
	if err != nil {
		goerrors.Log().Warnln("err on contest startTime Parse ", err)
		return false, err
	}
	startTimeUnix := startTime.Unix()
	now := time.Now().Unix()
	remained := now - startTimeUnix
	if remained >= 0 {
		return true, nil
	}
	return false, nil
}

func (c *Contest) BeforeDelete(tx *gorm.DB) (err error) {
	fmt.Println("BEFORE DELETE---------------------", c.ID, "---", c.StartTime)
	err = tx.Where("owner_id = ? and owner_type = ?", c.ID, "contests").Delete(&Media{}).Error
	if err != nil {
		return err
	}
	started, err := c.Started()
	if err != nil {
		goerrors.Log().Warnln("BEFORE DELETE check contest if started error: ", err)
		return err
	}
	if started {
		err = errors.New(fmt.Sprintf("Конкурс №%d уже начался", c.ID))
		goerrors.Log().Warnln("BEFORE DELETE ", err)
		return err
	}
	return
}

func (c *Contest) BeforeCreate(tx *gorm.DB) (err error) {
	//fmt.Println("BEFORE Create---------------------", c.ID)
	err = c.Validate()
	if err != nil {
		goerrors.Log().Warnln("BEFORECREATE contest validate error: ", err)
		return err
	}
	return
}
