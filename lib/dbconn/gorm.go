package dbconn

import (
	"gorm.io/driver/postgres"
	"gorm.io/plugin/prometheus"

	"gorm.io/gorm"
)

type Options struct {
	withoutPrometheus bool
}

type GormOptions func(o *Options)

func WithoutPrometheus() GormOptions {
	return func(o *Options) {
		o.withoutPrometheus = true
	}
}

func SetupGorm(dsn string, options ...GormOptions) (*gorm.DB, error) {
	o := &Options{}
	for _, f := range options {
		f(o)
	}
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: NewLogger(),
	})
	if err != nil {
		return nil, err
	}

	if o.withoutPrometheus {
		return gormDB, nil
	}

	err = gormDB.Use(prometheus.New(prometheus.Config{}))
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}
