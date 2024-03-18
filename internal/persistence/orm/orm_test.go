package orm

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestOrderLineEntityCanLoadLines(t *testing.T) {
	db := CreateInMemoryGormDb()
	AutoMigrate(db)

	db.Exec(`
	INSERT INTO order_lines (order_id, sku, qty) VALUES
        ('order1', 'RED-CHAIR', 12),
        ('order1', 'RED-TABLE', 13),
        ('order2', 'BLUE-LIPSTICK', 14)
	`)

	expected := []entity.OrderLine{
		{OrderID: "order1", SKU: "RED-CHAIR", Qty: 12},
		{OrderID: "order1", SKU: "RED-TABLE", Qty: 13},
		{OrderID: "order2", SKU: "BLUE-LIPSTICK", Qty: 14},
	}
	var got []entity.OrderLine
	db.Where(`(order_id = 'order1' and sku = 'RED-CHAIR') 
or (order_id = 'order1' and sku = 'RED-TABLE')
or (order_id = 'order2' and sku = 'BLUE-LIPSTICK')`).
		Select("order_id, sku, qty").Find(&got)
	assert.Equal(t, expected, got)
}

func TestOrderLaneEntityCanSaveLines(t *testing.T) {
	db := CreateInMemoryGormDb()
	AutoMigrate(db)

	expected := entity.OrderLine{
		OrderID: "order3",
		SKU:     "BLUE-LIPSTICK",
		Qty:     14,
	}
	db.Create(&expected)

	var got entity.OrderLine
	db.Raw("SELECT id, order_id, sku, qty FROM order_lines WHERE id = ?", expected.ID).Scan(&got)
	assert.Equal(t, expected, got)
}

func TestRetrievingBatches(t *testing.T) {
	db := CreateInMemoryGormDb()
	AutoMigrate(db)

	db.Exec(`
	INSERT INTO batches (reference, sku, purchased_quantity, eta)
    VALUES ('batch1', 'sku1', 100, null),
           ('batch2', 'sku1', 200, '2024-12-24')
	`)

	var batchID int64
	db.Raw("SELECT id FROM batches WHERE reference = ? and sku = ?", "batch1", "sku1").Scan(&batchID)
	var batchID2 int64
	db.Raw("SELECT id FROM batches WHERE reference = ? and sku = ?", "batch2", "sku1").Scan(&batchID2)

	db.Exec(`
	INSERT INTO order_lines (order_id, sku, qty) VALUES
       ('order1', 'sku1', 12)
	`)

	var lineID int64
	db.Raw("SELECT id FROM order_lines WHERE order_id = ? and sku = ?", "order1", "sku1").Scan(&lineID)

	db.Exec(`
	INSERT INTO allocations (order_line_id, batch_id) VALUES
       (?, ?)
	`, lineID, batchID)

	var allocationID int64
	db.Raw("SELECT id FROM allocations WHERE order_line_id = ? and batch_id = ?", lineID, batchID).Scan(&allocationID)

	line := entity.OrderLine{
		ID:      lineID,
		OrderID: "order1",
		SKU:     "sku1",
		Qty:     12,
	}
	batch := entity.Batch{
		ID:                batchID,
		Reference:         "batch1",
		SKU:               "sku1",
		PurchasedQuantity: 100,
		ETA:               time.Time{},
	}
	batch.Allocations = append(batch.Allocations, entity.Allocation{
		ID:          allocationID,
		OrderLine:   line,
		OrderLineID: lineID,
		BatchID:     batchID,
	})

	expected := []entity.Batch{
		batch,
		{
			ID:                batchID2,
			Reference:         "batch2",
			SKU:               "sku1",
			PurchasedQuantity: 200,
			ETA:               createDate(2024, 12, 24),
			Allocations:       make([]entity.Allocation, 0),
		},
	}

	fmt.Println([]int64{batchID, batchID2})
	var got []entity.Batch
	db.Model(&got).Preload("Allocations.OrderLine").Find(&got, []int64{batchID, batchID2})
	assert.Equal(t, expected, got)
	//assert.True(t, batchID > 0)
	//assert.True(t, lineID > 0)
	//assert.True(t, allocationID > 0)

}

func createDate(year, month, day int) time.Time {
	dateString := fmt.Sprintf("%v-%v-%v", year, month, day)
	date, error := time.Parse("2006-01-02", dateString)

	if error != nil {
		log.Println(error)
		panic(error)
		return time.Time{}
	}
	return date
}
