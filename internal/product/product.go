package product

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var validate = validator.New()

type Product struct {
	gorm.Model
	Name        string          `gorm:"type:varchar(255);not null" json:"name" validate:"required,max=255"`
	Description string          `gorm:"type:text;not null" json:"description" validate:"required"`
	Price       decimal.Decimal `gorm:"type:numeric(10,2);not null" json:"price"`
	Images      pq.StringArray  `gorm:"type:text[];not null;default:'{}'" json:"images" validate:"omitempty,dive,url"`
}

func New(name, description string, price decimal.Decimal, images []string) (*Product, error) {
	p := &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Images:      pq.StringArray(images),
	}
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Product) Validate() error {
	if p.Price.IsNegative() {
		return errors.New("price must not be negative")
	}
	return validate.Struct(p)
}

func (Product) TableName() string {
	return "products"
}
