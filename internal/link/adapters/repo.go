package adapters

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/ddd"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepo() *GormRepository {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Batches{})
	return &GormRepository{db: db}
}

func (repo *GormRepository) Add(batch ddd.Batch) {
	dbBatch := Batches{
		Sku:       batch.SKU,
		Reference: batch.Reference,
		Eta:       batch.ETA,
	}

	// Create
	repo.db.Create(&dbBatch)

}

func (repo *GormRepository) Get(reference string) ddd.Batch {
	return ddd.Batch{}
}

func (repo *GormRepository) List() []ddd.Batch {
	return make([]ddd.Batch, 0)
}
