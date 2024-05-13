package pkg

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiQuery struct {
	query  map[string]func(db *gorm.DB, value string) (*gorm.DB, error)
	values map[string]string
}

func NewApiQuery(query map[string]func(db *gorm.DB, value string) (*gorm.DB, error)) *ApiQuery {
	return &ApiQuery{query: query, values: map[string]string{}}
}

func (a *ApiQuery) SetQuery(c *gin.Context) {
	for key := range a.query {
		value := c.Query(key)
		if value != "" {
			a.values[key] = value
		}
	}
}

func (a *ApiQuery) ApplyQuery(db *gorm.DB) (*gorm.DB, error) {
	for key, value := range a.values {
		newDb, err := a.query[key](db, value)
		if err != nil {
			return nil, err
		}
		db = newDb
	}
	return db, nil
}
