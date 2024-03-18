package adapters

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/orm"
	"gorm.io/gorm"
)

type BatchRepo struct {
	db *gorm.DB
}

func NewBatchRepo() (*BatchRepo, error) {
	db := orm.CreateInMemoryGormDb()
	orm.AutoMigrate(db)
	return &BatchRepo{db: db}, nil
}

func (repo *BatchRepo) Add(ctx context.Context, batch *entity.Batch) {
	repo.db.WithContext(ctx).Create(batch)

}

func (repo *BatchRepo) Get(ctx context.Context, reference string) *entity.Batch {
	var batch entity.Batch
	repo.db.WithContext(ctx).Where("reference = ?", reference).
		Preload("Allocations.OrderLine").First(&batch)
	return &batch
}

func (repo *BatchRepo) List(ctx context.Context) []*entity.Batch {
	var batches []entity.Batch
	repo.db.WithContext(ctx).Preload("Allocations.OrderLine").Find(&batches)
	return toPointerEntityBatchList(batches)
}

//func (repo *BatchRepo) PointerList(ctx context.Context) []*entity.Batch {
//	batches := repo.List(ctx)
//	return toPointerEntityBatchList(batches)
//}
//

func toPointerEntityBatchList(batches []entity.Batch) []*entity.Batch {
	n := len(batches)
	batchesPtr := make([]*entity.Batch, n, n)
	for i, batch := range batches {
		batchesPtr[i] = &batch
	}
	return batchesPtr
}
