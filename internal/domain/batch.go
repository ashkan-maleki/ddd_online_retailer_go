package domain

import (
	"sort"
	"time"
)

type Batch struct {
	Reference         string
	SKU               string
	PurchasedQuantity int
	ETA               time.Time
	allocations       []OrderLine
}

func NewBatch(reference string, sku string, qty int, eta time.Time) *Batch {
	return &Batch{
		Reference:         reference,
		SKU:               sku,
		ETA:               eta,
		PurchasedQuantity: qty,
	}
}

func (b *Batch) Allocate(line OrderLine) {
	if b.allocations == nil {
		b.allocations = make([]OrderLine, 0)
	}
	if b.CanAllocate(line) {
		b.allocations = append(b.allocations, line)
	}
}

func (b *Batch) CanAllocate(line OrderLine) bool {
	return b.SKU == line.SKU && b.AvailableQuantity() >= line.Qty
}

func (b *Batch) AvailableQuantity() int {
	return b.PurchasedQuantity - b.AllocatedQuantity()
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
