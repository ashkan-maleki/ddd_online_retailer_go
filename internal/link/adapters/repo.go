package adapters

import (
	"context"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/orm"
	"gorm.io/gorm"
	"log"
)

type ProductRepo struct {
	db   *gorm.DB
	seen map[string]*entity.Product
}

func (repo *ProductRepo) Seen() map[string]*entity.Product {
	return repo.seen
}

func NewProductRepo() (*ProductRepo, error) {
	db := orm.CreateInMemoryGormDb()
	orm.AutoMigrate(db)
	return &ProductRepo{db: db, seen: make(map[string]*entity.Product)}, nil
}

func (repo *ProductRepo) Add(ctx context.Context, product *entity.Product) error {
	tx := repo.db.WithContext(ctx).Create(product)
	repo.saveSeenProduct(product)
	return tx.Error
}

func (repo *ProductRepo) Update(ctx context.Context, product *entity.Product) error {
	fmt.Println("bef update: ", product.Batches[0].PurchasedQuantity)
	tx := repo.db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Updates(product).Save(product)
	repo.saveSeenProduct(product)

	//fmt.Println("inside repo: ", len(product.Events()))
	return tx.Error
}

func (repo *ProductRepo) saveSeenProduct(product *entity.Product) {
	repo.seen[product.SKU] = product
}

func (repo *ProductRepo) Get(ctx context.Context, sku string) *entity.Product {
	var product entity.Product
	tx := repo.db.WithContext(ctx).Where("sku = ?", sku).
		Preload("Batches.Allocations.OrderLine").First(&product)
	if tx.Error != nil {
		return nil
	}
	repo.saveSeenProduct(&product)
	return &product
}

func (repo *ProductRepo) GetByBatchRef(ctx context.Context, ref string) *entity.Product {
	var product entity.Product
	tx := repo.db.WithContext(ctx).Preload("Batches", "reference = ?", ref).
		Preload("Batches.Allocations.OrderLine").First(&product)
	if tx.Error != nil {
		log.Println("GetByBatchRef: ", tx.Error)
		sql := repo.db.ToSQL(func(tx *gorm.DB) *gorm.DB {
			return tx.WithContext(ctx).Preload("Batches", "reference = ?", ref).
				Preload("Batches.Allocations.OrderLine").First(&product)
		})
		log.Println(sql)
		return nil
	}
	repo.saveSeenProduct(&product)
	return &product
}
