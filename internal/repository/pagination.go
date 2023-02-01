package repository

import (
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	Sort       string      `json:"sort,omitempty;query:sort"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Records    interface{} `json:"records"`
}

func Paginate(pagination *Pagination) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {

		offset := (pagination.Page - 1) * pagination.Limit
		return db.Offset(offset).Order(pagination.Sort).Limit(pagination.Limit)
	}
}

func GetPaginateSettings(r *http.Request) *Pagination {
	pagination := Pagination{}
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	if page == 0 {
		page = 1
	}
	pagination.Page = page

	limit, _ := strconv.Atoi(q.Get("limit"))
	switch {
	case limit > 50:
		limit = 50
	case limit <= 0:
		limit = 15
	}
	pagination.Limit = limit

	order := q.Get("order")
	if order == "" {
		order = "Id desc"
	}
	pagination.Sort = order

	return &pagination
}
