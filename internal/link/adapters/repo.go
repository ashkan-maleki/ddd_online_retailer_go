package adapters

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/orm"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db   *gorm.DB
	seen []*entity.Product
}

func (repo *ProductRepo) Seen() []*entity.Product {
	return repo.seen
}

func NewProductRepo() (*ProductRepo, error) {
	db := orm.CreateInMemoryGormDb()
	orm.AutoMigrate(db)
	return &ProductRepo{db: db}, nil
}

func (repo *ProductRepo) Add(ctx context.Context, product *entity.Product) error {
	tx := repo.db.WithContext(ctx).Create(product)
	repo.addSeenProduct(product)
	return tx.Error
}

func (repo *ProductRepo) Update(ctx context.Context, product *entity.Product) error {
	tx := repo.db.WithContext(ctx).Save(product)
	return tx.Error
}

func (repo *ProductRepo) addSeenProduct(product *entity.Product) {
	repo.seen = append(repo.seen, product)
}

func (repo *ProductRepo) Get(ctx context.Context, sku string) *entity.Product {
	var product entity.Product
	tx := repo.db.WithContext(ctx).Where("sku = ?", sku).
		Preload("Batches.Allocations.OrderLine").First(&product)
	if tx.Error != nil {
		return nil
	}
	repo.addSeenProduct(&product)
	return &product
}
