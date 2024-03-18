package orm

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/pkg/ddd/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

func CreateInMemoryGormDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(string(repository.InMemory)), &gorm.Config{})
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Batch{})
	if err != nil {
		log.Println(err)
		panic("migration failed")
	}
	err = db.AutoMigrate(&entity.OrderLine{})
	if err != nil {
		log.Println(err)
		panic("migration failed")
	}
	err = db.AutoMigrate(&entity.Allocation{})
	if err != nil {
		log.Println(err)
		panic("migration failed")
	}
}
