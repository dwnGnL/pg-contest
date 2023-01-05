package repository

func (r RepoImpl) CreateContest(contest Contest) (*Contest, error) {
	err := r.db.Create(&contest).Error
	if err != nil {
		return nil, err
	}
	return &contest, nil
}

func (r RepoImpl) DeleteContest(contest Contest) error {
	err := r.db.Select("Medias").Delete(&contest).Error
	if err != nil {
		return err
	}
	return nil
}

func (r RepoImpl) GetAllContest() (*[]Contest, error) {
	contest := new([]Contest)
	err := r.db.Preload("Questions.Answers").Preload("Medias").Find(&contest).Error
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (r RepoImpl) GetContest(contestID int64) (*Contest, error) {
	contest := new(Contest)
	err := r.db.Preload("Questions.Answers").Preload("Medias").Find(&contest, contestID).Error
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (r RepoImpl) UpdateContest(contest Contest) (*Contest, error) {
	err := r.db.Select("Medias").Delete(&contest).Error
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
