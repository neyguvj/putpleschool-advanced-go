package migrations

import (
	"fmt"

	"advancedgo/internal/product"

	"gorm.io/gorm"
)

func Migrate(gdb *gorm.DB) error {
	if err := gdb.AutoMigrate(
		&product.Product{},
	); err != nil {
		return fmt.Errorf("automigrate: %w", err)
	}
	return nil
}
