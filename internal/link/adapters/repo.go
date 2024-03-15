package adapters

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
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

func (repo *GormRepository) Add(batch domain.Batch) {
	dbBatch := Batches{
		Batch: batch,
	}

	// Create
	repo.db.Create(&dbBatch)

}

func (repo *GormRepository) Get(reference string) domain.Batch {
	return domain.Batch{}
}

func (repo *GormRepository) List() []domain.Batch {
	return make([]domain.Batch, 0)
}
