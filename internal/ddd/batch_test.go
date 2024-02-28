package ddd

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAllocatingToABatchReducesAvailableQuantity(t *testing.T) {
	batch := NewBatch("batch-001", "SMALL-TABLE", time.Now(), 20)
	line := NewOrderLine("order-ref", "SMALL-TABLE", 2)

	batch.Allocate(line)

	assert.Equal(t, 18, batch.AvailableQuantity())
}

func makeBatchAndLine(sku string, batchQty int, lineQty int) (*Batch, OrderLine) {
	return NewBatch("batch-001", sku, time.Now(), batchQty),
		NewOrderLine("order-123", sku, lineQty)
}

func TestCanAllocateIfAvailableGreaterThanRequired(t *testing.T) {
	largeBatch, smallLine := makeBatchAndLine("ELEGANT-LAMP", 20, 2)
	assert.True(t, largeBatch.CanAllocate(smallLine))
}

func TestCannotAllocateIfAvailableSmallerThanRequired(t *testing.T) {
	smallBatch, largeLine := makeBatchAndLine("ELEGANT-LAMP", 2, 20)
	assert.False(t, smallBatch.CanAllocate(largeLine))
}

func TestCanAllocateIfAvailableEqualToRequired(t *testing.T) {
	smallBatch, largeLine := makeBatchAndLine("ELEGANT-LAMP", 2, 2)
	assert.True(t, smallBatch.CanAllocate(largeLine))
}

func TestCannotAllocateIfSkusDoNotMatch(t *testing.T) {
	batch := NewBatch("batch-001", "UNCOMFORTABLE-CHAIR", time.Time{}, 100)
	differentSkuLine := NewOrderLine("order-123", "EXPENSIVE-TOASTER", 10)
	assert.False(t, batch.CanAllocate(differentSkuLine))
}

func TestCanOnlyDeallocateAllocatedLines(t *testing.T) {
	batch, unallocatedLine := makeBatchAndLine("DECORATIVE-TRINKET", 20, 2)
	batch.Deallocate(unallocatedLine)
	assert.Equal(t, 20, batch.AvailableQuantity())
}
