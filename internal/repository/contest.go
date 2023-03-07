package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"gorm.io/gorm/clause"
)

func (r RepoImpl) CreateContest(contest Contest) (*Contest, error) {
	err := r.db.Create(&contest).Error
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

func (r RepoImpl) DeleteContest(contest Contest) error {
	err := r.db.Select("Photos").Delete(&contest).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RepoImpl) GetAllContest(pagination *Pagination) (*Pagination, error) {

	var totalRows int64
	err := r.db.Model(Contest{}).Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	contest := new([]Contest)
	err = r.db.Scopes(Paginate(pagination)).Preload("Questions.Answers").Preload("Photos").Find(&contest).Error
	if err != nil {
		return nil, err
	}
	pagination.Records = contest
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(pagination.TotalRows / int64(pagination.Limit))
	if pagination.TotalRows%int64(pagination.Limit) > 0 {
		pagination.TotalPages++
	}
	return pagination, nil
}

func (r RepoImpl) GetAllContestByUserID(userID int64, pagination *Pagination) (*Pagination, error) {

	var totalRows int64
	err := r.db.Model(Contest{}).Where("NOT is_end and active").Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	userContestResp := new([]UserContestResp)
	//запрос для отображения всех конкурсов и с отображением наверху конкурсов которые купил данный участник (когда и за сколько)
	err = r.db.Table("contests c").
		Select("c.id AS id,"+
			"c.title AS title,"+
			"c.price AS price,"+
			"c.start_time AS start_time,"+
			"c.is_end AS is_end,"+
			//находим количество уникальных вопросов для каждого конкурса
			"COUNT(DISTINCT q.id) AS questions_count,"+
			//суммируем времена ответов на каждый из вопросов конкурса и дели на количесво фоток конкурса для устранения повторного суммирования, при делении обрабатываем случаё деления на 0
			"CAST(SUM(COALESCE(q.time, 0)) / COALESCE(NULLIF(COUNT(DISTINCT p.id), 0), 1) AS BIGINT) AS contest_length,"+
			//тут мы находим все линки ан фотки каждого конкурса, для случая когда фоток у конкурса нет то возвращаем {} инча Scan не сработает
			"CASE WHEN COUNT(DISTINCT p.id) = 0 THEN '{}' else ARRAY_AGG(DISTINCT p.link) END AS photos_links,"+
			"uc.created_at AS purchase_date,"+
			"uc.price AS purchase_price").
		Joins("LEFT OUTER JOIN questions q ON q.contest_id = c.id").
		Joins("LEFT OUTER JOIN photos p ON p.owner_id = c.id AND p.owner_type = ?", "contests").
		Joins("LEFT OUTER JOIN user_contests uc ON  uc.contest_id = c.id AND uc.user_id = ?", userID).
		Where("NOT c.is_end and c.active").
		Group("c.title, c.id, uc.created_at, c.price, c.start_time, c.is_end, uc.price").
		Order("uc.created_at ASC").Scopes(Paginate(pagination)).
		Scan(&userContestResp).Error
	if err != nil {
		return nil, err
	}
	pagination.Records = userContestResp
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(pagination.TotalRows / int64(pagination.Limit))
	if pagination.TotalRows%int64(pagination.Limit) > 0 {
		pagination.TotalPages++
	}
	return pagination, nil
}

func (r RepoImpl) GetContestFullStatsForUser(contestID, userID int64, currentQuestionOrder int) (contest *Contest, err error) {
	var (
		query            = r.db
		questionPosition int
		answerPosition   int
		i, j             int
		question         Question
		answer           Answer
	)
	if currentQuestionOrder > 0 {
		query = query.Preload("Questions", "order < ?", currentQuestionOrder)
	}
	query = query.Preload("Questions.Answers").Preload("Photos")

	if err = query.Find(&contest, contestID).Error; err != nil {
		return
	}

	var userAnswers []UserAnswers
	err = r.db.Table("user_answers").Where("contestID = ? AND userID = ?", contestID, userID).Scan(&userAnswers).Error
	if err != nil {
		return
	}
	for _, userAnswer := range userAnswers {
		questionPosition = -1
		answerPosition = -1
		for i, question = range contest.Questions {
			if question.ID == userAnswer.QuestionID {
				questionPosition = i
				for j, answer = range question.Answers {
					if answer.ID == userAnswer.AnswerID {
						answerPosition = j
						break
					}
				}
				break
			}
		}
		if questionPosition < 0 || answerPosition < 0 {
			err = errors.New(fmt.Sprintf("Check questionID = %v, answerID = %v availability in contestID = %v", userAnswer.QuestionID, userAnswer.AnswerID, contestID))
			return
		}
		contest.Questions[questionPosition].Answers[answerPosition].ChooseTime = userAnswer.Time
	}
	return
}

func (r RepoImpl) prepareContestStarQuery(contestID, currentQuestionID int64) *gorm.DB {
	/*return r.db.Table("user_contests uc").
	Select("row_number() over () AS rank,"+
		"uc.user_id AS user_id,"+
		"uc.user_name AS user_name,"+
		"COUNT(CASE is_correct WHEN true THEN 1 END ) AS total_correct,"+
		"SUM(CASE is_correct WHEN true THEN q.score ELSE 0 END) AS total_score,"+
		"SUM(CASE is_correct WHEN true THEN ua.time ELSE 0 END) AS total_time").
	Joins("LEFT OUTER JOIN user_answers ua ON uc.user_id = ua.user_id and ua.contest_id = uc.contest_id").
	Joins("LEFT OUTER JOIN answers a ON ua.answer_id = a.id AND ua.question_id = a.question_id AND ua.question_id <> ?", currentQuestionID).
	Joins("LEFT OUTER JOIN questions q ON q.id = ua.question_id").
	Where("uc.contest_id = ?", contestID).
	Group("uc.user_id, uc.user_name").
	Order("total_score DESC, total_time ASC")*/
	query := r.db.Table("user_contests uc").
		Select("uc.user_id AS user_id,"+
			"uc.user_name AS user_name,"+
			"COUNT(CASE is_correct WHEN true THEN 1 END ) AS total_correct,"+
			"SUM(CASE is_correct WHEN true THEN q.score ELSE 0 END) AS total_score,"+
			"SUM(CASE is_correct WHEN true THEN ua.time ELSE 0 END) AS total_time").
		Joins("LEFT OUTER JOIN user_answers ua ON uc.user_id = ua.user_id and ua.contest_id = uc.contest_id").
		Joins("LEFT OUTER JOIN answers a ON ua.answer_id = a.id AND ua.question_id = a.question_id AND ua.question_id <> ?", currentQuestionID).
		Joins("LEFT OUTER JOIN questions q ON q.id = ua.question_id").
		Where("uc.contest_id = ?", contestID).
		Group("uc.user_id, uc.user_name").
		Order("total_score DESC, total_time ASC")
	return r.db.Table("(?) as a", query).Select("row_number() over () AS rank, a.*")
}

func (r RepoImpl) GetContestStatsForUser(contestID, userID, currentQuestionID int64) (contestStatsResp *ContestStats, err error) {
	err = r.db.Table("(?) as b", r.prepareContestStarQuery(contestID, currentQuestionID)).Where("b.user_id = ?", userID).Scan(&contestStatsResp).Error
	return
}

func (r RepoImpl) GetContestStatsById(contestID, currentQuestionID int64, pagination *Pagination) (*Pagination, error) {

	var totalRows int64
	//считаем общее количество пользователей купивших данный конкурс
	err := r.db.Model(UserContests{}).Where("contest_id = ?", contestID).Count(&totalRows).Error
	if err != nil {
		return nil, err
	}

	contestStatsResp := new([]ContestStats)
	//запрос для отображения результатов ВСЕХ пользователей купивших данный конкурс, независимо от факта участия
	err = r.prepareContestStarQuery(contestID, currentQuestionID).Scopes(Paginate(pagination)).Scan(&contestStatsResp).Error
	if err != nil {
		return nil, err
	}

	pagination.Records = contestStatsResp
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(pagination.TotalRows / int64(pagination.Limit))
	if pagination.TotalRows%int64(pagination.Limit) > 0 {
		pagination.TotalPages++
	}
	return pagination, nil
}

func (r RepoImpl) GetContest(contestID int64) (contest *Contest, err error) {
	err = r.db.Preload("Questions.Answers").Preload("Photos").Last(&contest, contestID).Error
	return
}

func (r RepoImpl) GetContestInfo(contestID int64) (contest *Contest, err error) {
	err = r.db.Table("contests").Where("id = ?", contestID).Last(&contest).Error
	return
}

func (r RepoImpl) UpdateContest(contest Contest) (*Contest, error) {
	err := r.db.Select("Photos").Delete(&contest).Error
	if err != nil {
		return nil, err
	}
	err = r.db.Create(&contest).Error
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

func (r RepoImpl) ChangeContestInfo(contest *Contest) error {
	return r.db.Updates(&contest).Error
}

func (r RepoImpl) SubscribeContest(userContest *UserContests) error {
	return r.db.Create(&userContest).Error
}

func (r RepoImpl) ContestAvailability(contestID int64, userID int64) (contest *Contest, err error) {

	var userContests int64
	err = r.db.Table("user_contests").Where("user_id = ? AND contest_id = ?", userID, contestID).Count(&userContests).Error
	if err != nil {
		return
	}
	if userContests != 0 {
		err = errors.New("already subscribed earlier")
		return
	}

	err = r.db.Table("contests").Where("active AND NOT is_end AND id = ?", contestID).Last(&contest).Error
	return
}

func (r RepoImpl) GetUserContest(contestID int64, userID int64) (userContest *UserContests, err error) {
	err = r.db.Where("user_id = ? and contest_id = ?", userID, contestID).Find(&userContest).Error
	return
}

func (r RepoImpl) SubmitAnswer(userAnswer *UserAnswers) (err error) {
	var count int64

	err = r.db.Model(UserAnswers{}).Where("user_id = ? and contest_id = ? and question_id = ? and answer_id = ?",
		userAnswer.UserID, userAnswer.ContestID, userAnswer.QuestionID, userAnswer.AnswerID).
		Count(&count).Error
	if err != nil {
		return
	}
	//усли участник уже так ответил ранееб не обновляем ничего
	if count != 0 {
		return
	}

	err = r.db.Clauses(clause.OnConflict{
		//Columns:   []clause.Column{{Name: "user_id"}, {Name: "question_id"}, {Name: "contest_id"}},
		UpdateAll: true,
	}).Create(userAnswer).Error
	return
}
