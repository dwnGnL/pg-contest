package repository

func (r RepoImpl) GetUserTikets(userID, tiketID int64) (*UserTickets, error) {
	tikets := new(UserTickets)
	err := r.db.Where("user_id = ? and contest_id = ?", userID, tiketID).Find(tikets).Error
	if err != nil {
		return nil, err
	}
	return tikets, nil
}
