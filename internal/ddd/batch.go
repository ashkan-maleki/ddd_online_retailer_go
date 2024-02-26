package ddd

import (
	"errors"
	"time"
)

var OutOfStock = errors.New("out of stock")

type OrderLine struct {
	OrderId string
	SKU     string
	Qty     int
}

func NewOrderLine(orderId string, SKU string, qty int) OrderLine {
	return OrderLine{OrderId: orderId, SKU: SKU, Qty: qty}
}

func (ol OrderLine) EqualTo(line OrderLine) bool {
	return ol.OrderId == line.OrderId && ol.Qty == line.Qty && ol.SKU == line.SKU
}

type Batch struct {
	Reference         string
	SKU               string
	ETA               time.Time
	purchasedQuantity int
	allocations       []OrderLine
}

func (b *Batch) Allocate(line OrderLine) {
	if b.CanAllocate(line) {
		b.allocations = append(b.allocations, line)
	}
}

func (b *Batch) CanAllocate(line OrderLine) bool {
	return b.SKU == line.SKU && b.AvailableQuantity() >= line.Qty
}

func (b *Batch) AvailableQuantity() int {
	return b.purchasedQuantity - b.AllocatedQuantity()
}

func (b *Batch) AllocatedQuantity() int {
	sum := 0
	for _, line := range b.allocations {
		sum += line.Qty
	}
	return sum
}

func NewBatch(reference string, sku string, eta time.Time, qty int) *Batch {
	return &Batch{
		Reference:         reference,
		SKU:               sku,
		ETA:               eta,
		purchasedQuantity: qty,
		allocations:       make([]OrderLine, 0),
	}
}
