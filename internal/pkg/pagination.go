package pkg

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      int         `json:"limit,omitempty;query:limit"`
	Page       int         `json:"page,omitempty;query:page"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func (p *Pagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (pagination *Pagination) Paginate(db *gorm.DB) *gorm.DB {
	return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit())

}

func (pagination *Pagination) Count(db *gorm.DB) {
	var totalRows int64
	db.Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
}

func (pagination *Pagination) CountRaw(db *gorm.DB, query string) {
	var totalRows int64
	db.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM (%s)`, query)).Scan(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages
}
