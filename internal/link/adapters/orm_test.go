package adapters

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/pkg/ddd/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"testing"
	"time"
)

func TestOrderLineEntityCanLoadLines(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(string(repository.InMemory)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&OrderLines{})
	if err != nil {
		return
	}

	db.Exec(`
	INSERT INTO order_lines (order_id, sku, qty) VALUES
        ('order1', 'RED-CHAIR', 12),
        ('order1', 'RED-TABLE', 13),
        ('order2', 'BLUE-LIPSTICK', 14)
	`)

	expected := []OrderLines{
		{OrderLine: domain.OrderLine{OrderID: "order1", SKU: "RED-CHAIR", Qty: 12}},
		{OrderLine: domain.OrderLine{OrderID: "order1", SKU: "RED-TABLE", Qty: 13}},
		{OrderLine: domain.OrderLine{OrderID: "order2", SKU: "BLUE-LIPSTICK", Qty: 14}},
	}
	var got []OrderLines
	db.Find(&got)
	assert.Equal(t, MapManyDomainOrderLine(expected), MapManyDomainOrderLine(got))
}

func TestOrderLaneEntityCanSaveLines(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(string(repository.InMemory)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&OrderLines{})

	// Create
	expected := OrderLines{OrderLine: domain.OrderLine{OrderID: "order2",
		SKU: "BLUE-LIPSTICK", Qty: 14}}
	db.Create(&expected)

	var got OrderLines
	db.Raw("SELECT id, order_id, sku, qty FROM order_lines WHERE id = ?", expected.ID).Scan(&got)
	assert.Equal(t, expected, got)
}

func GetSchema(table any, db *gorm.DB) *schema.Schema {
	stmt := &gorm.Statement{DB: db}
	stmt.Parse(table)
	return stmt.Schema
}

func TestRetrievingBatches(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(string(repository.InMemory)), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	//fmt.Println("8888888888888888")
	err = db.AutoMigrate(&Batches{})
	err = db.AutoMigrate(&OrderLines{})
	err = db.AutoMigrate(&Allocations{})
	if err != nil {
		log.Println(err)
	}

	//fmt.Println(GetSchema(&Batches{}, db))

	db.Exec(`
	 INSERT INTO batches (reference, sku, purchased_quantity, eta)
     VALUES ('batch1', 'sku1', 100, null)
	`)

	db.Exec(`
	 INSERT INTO batches (reference, sku, purchased_quantity, eta)
     VALUES ('batch2', 'sku1', 200, '2024-12-24')
	`)

	db.Exec(`
	INSERT INTO order_lines (order_id, sku, qty) VALUES
        ('order1', 'RED-CHAIR', 12)
	`)

	//line := OrderLines{OrderLine: domain.OrderLine{OrderID: "order1", SKU: "sku1", Qty: 12}}

	db.Exec(`
	INSERT INTO allocations (order_line_id, batch_id) VALUES
        (1,1)
	`)

	expected := []Batches{
		{Batch: *domain.NewBatch("batch1", "sku1", 100, time.Time{})},
		{Batch: *domain.NewBatch("batch2", "sku1", 200,
			createDate(2024, 12, 24))},
	}

	expected[0].Allocations = append(expected[0].Allocations, Allocations{BatchID: 1, OrderLineID: 1})

	var got []Batches
	db.Model(&got).Preload("Allocations").Find(&got)
	assert.Equal(t, MapManyDomainBatch(expected), MapManyDomainBatch(got))
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
