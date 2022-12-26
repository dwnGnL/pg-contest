package models

type ContestStatus int

const (
	Waiting ContestStatus = iota + 1
	Start
	End
)

type WsRequest struct {
	Token      string `json:"token"`
	QuestionID int64  `json:"question_id"`
	AnswerID   int64  `json:"answer_id"`
}

type WsResponse struct {
	Step             int           `json:"step"`
	TotalStep        int           `json:"total_step"`
	ContestStatus    ContestStatus `json:"contest_status"`
	ActiveQuestionID int64         `json:"active_question_id"`
	CountDown        int64         `json:"count_down"`
	TotalTime        int64         `json:"total_time"`
	Questions        []WsQuestion  `json:"questions"`
}

type WsQuestion struct {
	ID      int64      `json:"id"`
	Order   int        `json:"order"`
	Title   string     `json:"title"`
	Answers []WsAnswer `json:"answers"`
}

type WsAnswer struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}
