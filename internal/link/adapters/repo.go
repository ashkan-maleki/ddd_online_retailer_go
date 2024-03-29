package adapters

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/orm"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewBatchRepo() (*ProductRepo, error) {
	db := orm.CreateInMemoryGormDb()
	orm.AutoMigrate(db)
	return &ProductRepo{db: db}, nil
}

func (repo *ProductRepo) Add(ctx context.Context, product *entity.Product) {
	repo.db.WithContext(ctx).Create(product)

}

func (repo *ProductRepo) Get(ctx context.Context, sku string) *entity.Product {
	var product entity.Product
	repo.db.WithContext(ctx).Where("sku = ?", sku).
		Preload("Batches.Allocations.OrderLine").First(&product)
	return &product
}
