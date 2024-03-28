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
	Product           Product `gorm:"foreignKey:SKU;references:SKU"`
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

type Product struct {
	SKU           string  `gorm:"primaryKey"`
	VersionNumber int     `gorm:"default:0"`
	Batches       []Batch `gorm:"foreignKey:SKU;references:SKU"`
}
