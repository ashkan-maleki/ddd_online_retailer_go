package entity

import (
	"time"
)

type OrderLine struct {
	ID      int64 `gorm:"primaryKey"`
	OrderID string
	SKU     string
	Qty     int
	//Allocations []Allocation `gorm:"foreignKey:OrderLineID"`
}

type Batch struct {
	ID                int64 `gorm:"primaryKey"`
	Reference         string
	SKU               string
	PurchasedQuantity int
	ETA               time.Time
	Allocations       []Allocation `gorm:"foreignKey:BatchID"`
}

type Allocation struct {
	ID          int64 `gorm:"primaryKey"`
	OrderLineID int64
	OrderLine   OrderLine `gorm:"foreignKey:ID;references:OrderLineID"`
	BatchID     int64
}
