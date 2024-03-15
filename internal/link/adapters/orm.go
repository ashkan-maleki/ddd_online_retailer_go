package adapters

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
)

type OrderLines struct {
	ID int64 `gorm:"primaryKey"`
	domain.OrderLine
	//Sku string
	//Qty     int
	//OrderID string
}

func GetDomainOrderLine(lines []OrderLines) []domain.OrderLine {
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

type Allocations struct {
	ID          int64 `gorm:"primaryKey"`
	OrderLineID int64
	BatchID     int64
}
