package presentation

import "time"

type Batch struct {
	Reference         string
	SKU               string
	PurchasedQuantity int
	ETA               time.Time
}

type OrderLine struct {
	OrderID string
	SKU     string
	Qty     int
}
