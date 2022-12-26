package repository

func (r RepoImpl) GetContest(contestID int64) (*Contest, error) {
	contest := new(Contest)
	err := r.db.Preload("Questions.Answers").Find(contest, contestID).Error
	if err != nil {
		return nil, err
	}
	return contest, nil
}

func (r RepoImpl) UpdateContest(contestID int64, contest Contest) error {
	err := r.db.Where("id", contestID).Updates(contest).Error
	if err != nil {
		return err
	}
	return nil
}
