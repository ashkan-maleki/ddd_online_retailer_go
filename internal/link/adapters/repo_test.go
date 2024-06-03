package adapters

import (
	"context"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/orm"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

func TestRepositoryCanSaveABatchAndAProduct(t *testing.T) {

	repo := newBatchRepo()
	sku := "RUSTY-SOAPDISH"
	batch := entity.Batch{
		Reference:         "batch5",
		SKU:               sku,
		PurchasedQuantity: 100,
		ETA:               time.Time{},
	}
	expected := entity.Product{
		SKU:           sku,
		VersionNumber: 0,
		Batches:       []entity.Batch{batch},
	}
	repo.Add(context.Background(), &expected)

	db := orm.CreateInMemoryGormDb()

	var got entity.Product
	db.Where("sku = ?", batch.SKU).Preload("Batches").First(&got)
	assert.Equal(t, expected, got)
}

func newBatchRepo() *ProductRepo {
	repo, err := NewProductRepo()
	if err != nil {
		log.Println(err)
		panic("new batch repo failed")
	}
	return repo
}

func InsertProduct(db *gorm.DB) {
	product := entity.Product{
		SKU:           "GENERIC-SOFA",
		VersionNumber: 0,
		Batches:       make([]entity.Batch, 0),
	}
	db.Create(&product)
}

func InsertOrderLine(db *gorm.DB) int64 {
	line := entity.OrderLine{
		OrderID: "order3",
		SKU:     "GENERIC-SOFA",
		Qty:     12,
	}
	db.Create(&line)
	return line.ID
	//db.Exec(`
	//INSERT INTO order_lines (order_id, sku, qty) VALUES
	//    ('order3', 'GENERIC-SOFA', 12)
	//`)
	//
	//var got entity.OrderLine
	//db.Raw("SELECT id FROM order_lines WHERE order_id = ? and sku = ?", "order3", "GENERIC-SOFA").Scan(&got)
	//return got.TransactionID
}

func InsertBatch(db *gorm.DB, reference string) int64 {
	batch := entity.Batch{
		Reference:         reference,
		SKU:               "GENERIC-SOFA",
		PurchasedQuantity: 100,
		ETA:               time.Time{},
	}
	db.Create(&batch)
	return batch.ID

	//db.Exec(`
	// INSERT INTO batches (reference, sku, purchased_quantity, eta)
	// VALUES (?, 'GENERIC-SOFA', 100, null)
	//`, reference)
	//
	//var got entity.Batch
	//db.Raw("SELECT id FROM batches WHERE reference = ? and sku = ?", reference, "GENERIC-SOFA").Scan(&got)
	//return got.TransactionID
}

func InsertAllocation(db *gorm.DB, lineID int64, batchID int64) {
	allocation := entity.Allocation{
		OrderLineID: lineID,
		BatchID:     batchID,
	}
	db.Create(&allocation)
	//db.Exec(`
	// INSERT INTO allocations (order_line_id, batch_id)
	//VALUES (?, ?)
	//`, lineID, batchID)
	//
	//var alocs []entity.Allocation
	//db.Raw("select * from allocations").Find(&alocs)
	//log.Println("======")
	//log.Println(alocs)
	//log.Println("======")
}

func TestRepositoryCanRetrieveABatchWithAllocations(t *testing.T) {
	repo := newBatchRepo()
	db := orm.CreateInMemoryGormDb()
	InsertProduct(db)
	lineID := InsertOrderLine(db)
	batchID := InsertBatch(db, "batch3")
	InsertBatch(db, "batch4")
	InsertAllocation(db, lineID, batchID)

	var expected entity.Batch
	db.Preload("Allocations.OrderLine").First(&expected, batchID)

	fmt.Println("=============")
	fmt.Println(expected)
	fmt.Println("=============")
	batch := repo.Get(context.Background(), "GENERIC-SOFA").Batches[0]
	fmt.Println(batch)
	got := batch
	assert.Equal(t, expected, got)
	assert.Equal(t, 1, len(batch.Allocations))

}
