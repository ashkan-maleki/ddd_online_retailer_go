package ddd

import (
	"testing"
	"time"
)

func TestAllocatingToABatchReducesAvailableQuantity(t *testing.T) {
	batch := NewBatch("batch-001", "SMALL-TABLE", time.Now(), 20)
	line := NewOrderLine("order-ref", "SMALL-TABLE", 2)

	batch.Allocate(line)
	t.
}
