package product

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrNotFound = errors.New("product not found")

type Repository interface {
	Create(ctx context.Context, p *Product) error
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*Product, error)
}

type gormRepository struct {
	gdb *gorm.DB
}

func NewRepository(gdb *gorm.DB) Repository {
	return &gormRepository{gdb: gdb}
}

func (r *gormRepository) Create(ctx context.Context, p *Product) error {
	if err := r.gdb.WithContext(ctx).Create(p).Error; err != nil {
		return err
	}
	return nil
}

func (r *gormRepository) Update(ctx context.Context, p *Product) error {
	res := r.gdb.WithContext(ctx).
		Model(&Product{}).
		Where("id = ?", p.ID).
		Select("name", "description", "price", "images").
		Updates(p)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepository) Delete(ctx context.Context, id uint) error {
	res := r.gdb.WithContext(ctx).Delete(&Product{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *gormRepository) GetByID(ctx context.Context, id uint) (*Product, error) {
	var p Product
	if err := r.gdb.WithContext(ctx).First(&p, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}
