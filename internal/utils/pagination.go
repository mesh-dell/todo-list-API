package utils

import (
	"context"
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Data       any `json:"data"`
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	TotalPages int `json:"total"`
}

func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 {
		p.Limit = 10
	}
}

func Paginate[T any](db *gorm.DB, ctx context.Context, p *Pagination, out *[]T) (*Pagination, error) {
	p.Normalize()

	var total int64
	db.Model(new(T)).Count(&total)

	p.TotalPages = int(math.Ceil(float64(total) / float64(p.Limit)))
	offset := (p.Page - 1) * p.Limit

	err := db.Limit(p.Limit).Offset(offset).Find(out).Error
	if err != nil {
		return nil, err
	}
	p.Data = out
	return p, nil
}
