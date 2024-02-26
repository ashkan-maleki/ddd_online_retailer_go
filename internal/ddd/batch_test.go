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
