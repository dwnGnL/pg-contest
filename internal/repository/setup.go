package repository

import (
	"github.com/dwnGnL/pg-contests/internal/config"
	"github.com/dwnGnL/pg-contests/lib/dbconn"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewRepository(cfg *config.Config) (*RepoImpl, error) {
	gormDB, err := dbconn.SetupGorm(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	return &RepoImpl{
		db: gormDB,
	}, nil
}

func (r RepoImpl) Migrate() error {
	for _, model := range []interface{}{
		(*Contest)(nil),
		(*Question)(nil),
		(*Answer)(nil),
		(*UserTickets)(nil),
	} {
		dbSilent := r.db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)})
		err := dbSilent.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

type RepoImpl struct {
	db *gorm.DB
}
