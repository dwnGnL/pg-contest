package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/dwnGnL/pg-contests/lib/goerrors"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

var layout = "2006-01-02T15:04Z07:00"

type UserContestResp struct {
	ID             int64          `json:"id" gorm:"column:id"`
	Title          string         `json:"title" gorm:"column:title"`
	Price          float64        `json:"price" gorm:"column:price"`
	StartTime      string         `json:"start_time" gorm:"column:start_time"`
	PhotosLinks    pq.StringArray `json:"photos_links" gorm:"column:photos_links"`
	QuestionsCount int64          `json:"questions_count" gorm:"column:questions_count"`
	ContestLength  int64          `json:"contest_length" gorm:"column:contest_length"`
	IsEnd          *bool          `json:"is_end" gorm:"column:is_end"`
	PurchaseDate   *time.Time     `json:"purchase_date" gorm:"column:purchase_date"`
	PurchasePrice  *float64       `json:"purchase_price" gorm:"column:purchase_price"`
}

type Contest struct {
	ID           int64      `json:"id" gorm:"column:id;primary_key;autoIncrement"`
	Title        string     `json:"title" binding:"required" gorm:"column:title"`
	Price        float64    `json:"price" binding:"required" gorm:"column:price"`
	PlayersCount *int64     `json:"players_count" gorm:"players_count"`
	StartTime    string     `json:"start_time" binding:"required" gorm:"column:start_time"`
	CreatedBy    string     `json:"created_by" gorm:"column:created_by"`
	Photos       []Photo    `json:"photos" gorm:"polymorphic:Owner;constraint:OnDelete:CASCADE;"`
	Questions    []Question `json:"questions" gorm:"foreignKey:ContestID;constraint:OnDelete:CASCADE"`
	Active       *bool      `json:"active" gorm:"column:active;default:true"`
	IsEnd        *bool      `json:"is_end" gorm:"column:is_end;default:false"`
	CreatedAt    *time.Time `json:"created_at" gorm:"autoCreateTime"`
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
	QuestionID int64  `json:"question_id" binding:"required" gorm:"column:question_id"`
	Title      string `json:"title" binding:"required" gorm:"column:title"`
	IsCorrect  *bool  `json:"is_correct" gorm:"column:is_correct;default:false"`
}

type Photo struct {
	ID        int64  `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	FileName  string `json:"file_name,omitempty" binding:"required" gorm:"column:file_name"`
	Uploaded  *bool  `json:"uploaded,omitempty" gorm:"column:uploaded;default:false"`
	Link      string `json:"link,omitempty" gorm:"column:link"`
	OwnerID   int64  `json:"owner_id,omitempty" gorm:"column:owner_id"`
	OwnerType string `json:"owner_type,omitempty" gorm:"column:owner_type"`
}

type UserContests struct {
	UserID    int64      `json:"user_id" gorm:"column:user_id;primaryKey"`
	ContestID int64      `json:"contest_id" gorm:"column:contest_id;primaryKey"`
	Price     float64    `json:"price" gorm:"column:price"`
	CreatedAt *time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type UserAnswers struct {
	UserID     int64 `json:"user_id" gorm:"column:user_id;primaryKey"`
	ContestID  int64 `json:"contest_id" gorm:"column:contest_id;primaryKey"`
	QuestionID int64 `json:"question_id" gorm:"column:question_id;primaryKey"`
	AnswerID   int64 `json:"answer_id" gorm:"column:answer_id"`
	Time       int64 `json:"time" gorm:"column:time"`
}

type UserTickets struct {
	UserID    int64 `gorm:"column:user_id"`
	ContestID int64 `gorm:"column:contest_id"`
	Canseled  bool  `gorm:"column:canseled;default:false"`
}

type ErrorResponse struct {
	Error ErrorStruct `json:"error"`
}

type ErrorStruct struct {
	Message string `json:"message,omitempty"`
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
			if *answer.IsCorrect {
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
	err = tx.Where("owner_id = ? and owner_type = ?", c.ID, "contests").Delete(&Photo{}).Error
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

type NullString64Array struct {
	StringArray pq.StringArray
	Valid       bool
}

func (n *NullString64Array) Scan(value interface{}) error {
	if value == nil {
		n.StringArray, n.Valid = nil, false
		return nil
	}
	n.Valid = true
	return pq.Array(&n.StringArray).Scan(value)
}
