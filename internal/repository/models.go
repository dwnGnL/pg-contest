package repository

type Contest struct {
	ID        int64      `gorm:"column:id;primary_key;autoIncrement"`
	Title     string     `gorm:"column:title"`
	Price     float64    `gorm:"column:price"`
	StartTime int64      `gorm:"column:start_time"`
	CreatedBy string     `gorm:"column:created_by"`
	Questions []Question `gorm:"foreignKey:ContestID"`
	Active    bool       `gorm:"column:active;default:false"`
	IsEnd     bool       `gorm:"column:is_end;default:false"`
}

type Question struct {
	ID        int64    `gorm:"column:id;primaryKey;autoIncrement"`
	ContestID int64    `gorm:"column:contest_id"`
	Answers   []Answer `gorm:"foreignKey:QuestionID"`
	Title     string   `gorm:"column:title"`
	Order     int      `gorm:"column:sort_order"`
	Time      int64    `gorm:"column:time"`
}

type Answer struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement"`
	QuestionID int64  `gorm:"column:contest_id"`
	Title      string `gorm:"column:title"`
	IsCorrect  bool   `gorm:"column:is_correct;default:false"`
}

type UserTickets struct {
	UserID    int64 `gorm:"column:user_id"`
	ContestID int64 `gorm:"column:contest_id"`
	Canseled  bool  `gorm:"column:canseled;default:false"`
}
