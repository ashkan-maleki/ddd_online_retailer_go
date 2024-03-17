package adapters

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/pkg/ddd/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type BatchRepo struct {
	db *gorm.DB
}

func NewBatchRepo() (*BatchRepo, error) {
	db, err := gorm.Open(sqlite.Open(string(repository.InMemory)), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Batches{})
	if err != nil {
		return nil, err
	}
	return &BatchRepo{db: db}, nil
}

func (repo *BatchRepo) Add(ctx context.Context, batch domain.Batch) {
	dbBatch := Batches{
		Batch: batch,
	}
	repo.db.WithContext(ctx).Create(&dbBatch)

}

func (repo *BatchRepo) Get(ctx context.Context, reference string) domain.Batch {
	var batch *Batches
	repo.db.WithContext(ctx).Where("reference = ?", reference).First(&batch)
	return batch.Batch
}

func (repo *BatchRepo) List(ctx context.Context) []domain.Batch {
	var batchList []Batches
	repo.db.WithContext(ctx).Find(&batchList)
	return MapManyDomainBatch(batchList)
}
