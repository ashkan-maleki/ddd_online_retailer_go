package domain

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

var OutOfStockErr = errors.New("out of stock")

type Batch struct {
	ID                int64
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
		PurchasedQuantity: qty,
		ETA:               eta,
		allocations:       make([]OrderLine, 0),
	}
}

func (b *Batch) Allocate(line OrderLine) {
	if b.CanAllocate(line) {
		fmt.Printf("Dom 1 --> OrderID: %v, allocation size: %v\n", line.OrderID, len(b.allocations))
		b.allocations = append(b.allocations, line)
		fmt.Printf("Dom 2 --> OrderID: %v, allocation size: %v\n", line.OrderID, len(b.allocations))
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

func (b *Batch) Allocations() []OrderLine {
	return b.allocations
}

func (b *Batch) DeallocateOne() *OrderLine {
	if len(b.allocations) > 1 {
		b.allocations = b.allocations[1:]
		return &b.allocations[0]
	} else if len(b.allocations) == 1 {
		b.allocations = make([]OrderLine, 0)
		return &b.allocations[0]
	} else {
		return nil
	}
}

// Allocate TODO: remove this function
// Allocate deprecate
func Allocate(line OrderLine, batches []*Batch) (*Batch, error) {
	sort.Slice(batches, func(i, j int) bool {
		return batches[i].ETA.Before(batches[j].ETA)
	})
	for _, batch := range batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			return batch, nil
		}
	}
	return nil, OutOfStockErr
}
