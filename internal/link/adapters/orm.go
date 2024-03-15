package adapters

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
)

type OrderLines struct {
	ID int64 `gorm:"primaryKey"`
	domain.OrderLine
	Allocations []Allocations `gorm:"foreignKey:OrderLineID"`
}

func MapManyDomainOrderLine(lines []OrderLines) []domain.OrderLine {
	domainLines := make([]domain.OrderLine, len(lines), len(lines))
	for i, line := range lines {
		domainLines[i] = line.OrderLine
	}
	return domainLines
}

type Batches struct {
	ID int64 `gorm:"primaryKey"`
	domain.Batch
	Allocations []Allocations `gorm:"foreignKey:BatchID"`
}

func MapManyDomainBatch(batches []Batches) []domain.Batch {
	domainBatches := make([]domain.Batch, len(batches), len(batches))
	for i, batch := range batches {
		domainBatches[i] = MapDomainBatch(batch)
	}
	return domainBatches
}

func MapDomainBatch(batch Batches) domain.Batch {
	b := batch.Batch
	fmt.Println("**** 1")
	for _, allocation := range batch.Allocations {
		b.Allocate(allocation.OrderLine.OrderLine)
	}
	return b
}

type Allocations struct {
	ID          int64 `gorm:"primaryKey"`
	OrderLineID int64
	OrderLine   OrderLines
	BatchID     int64
	Batch       Batches
}
