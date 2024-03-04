package ddd

import (
	"errors"
	"sort"
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

func (b *Batch) EqualTo(other Batch) bool {
	return b.Reference == other.Reference
}

func (b *Batch) Contain(line OrderLine) (bool, int) {
	for i, ol := range b.allocations {
		if ol.EqualTo(line) {
			return true, i
		}
	}
	return false, -1
}

func (b *Batch) Deallocate(line OrderLine) {
	ok, idx := b.Contain(line)
	if ok {
		b.allocations = append(b.allocations[:idx], b.allocations[idx+1:]...)
	}
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

func Allocate(line OrderLine, batches []*Batch) (string, error) {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].ETA.Before(batches[j].ETA)
	})
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			return batch.Reference, nil
		}
	}
	return "", OutOfStock

}
