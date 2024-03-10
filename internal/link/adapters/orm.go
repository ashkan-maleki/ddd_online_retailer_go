package adapters

import (
	"time"
)

type OrderLines struct {
	Id      int64 `gorm:"primaryKey"`
	Sku     string
	Qty     int
	OrderId int64
}

type Batches struct {
	Id          int64 `gorm:"primaryKey"`
	Reference   string
	Sku         string
	Eta         time.Time
	Allocations []Allocations
}

type Allocations struct {
	Id          int64 `gorm:"primaryKey"`
	OrderLineId int64
	BatchId     int64
}
