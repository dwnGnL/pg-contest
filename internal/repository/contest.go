package repository

import "errors"

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

func (r RepoImpl) GetContest(contestID int64) (contest *Contest, err error) {
	err = r.db.Preload("Questions.Answers").Preload("Photos").Last(&contest, contestID).Error
	return
}

func (r RepoImpl) GetContestPrice(contestID int64) (contest *Contest, err error) {
	err = r.db.Table("contests c").Select("id, price").Where("id = ?", contestID).Scan(&contest).Error
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

func (r RepoImpl) ChangeContestInfo(contest Contest) (*Contest, error) {
	err := r.db.Updates(&contest).Error
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

func (r RepoImpl) SubscribeContest(contest Contest, userID int64) (err error) {
	userContest := UserContests{
		UserID:    userID,
		ContestID: contest.ID,
		Price:     contest.Price,
	}
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
