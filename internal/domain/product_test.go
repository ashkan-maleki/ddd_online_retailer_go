package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func now() time.Time {
	return time.Now()
}

func tomorrow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

func later() time.Time {
	return time.Now().Add(10 * 24 * time.Hour)
}

func TestPrefersCurrentStockBatchesToShipments(t *testing.T) {
	inStockBatch := NewBatch("in-stock-batch", "RETRO-CLOCK", 100, time.Time{})
	shipmentBatch := NewBatch("shipment-batch", "RETRO-CLOCK", 100, tomorrow())
	line := NewOrderLine("oref", "RETRO-CLOCK", 10)

	_, err := Allocate(line, []*Batch{inStockBatch, shipmentBatch})
	if err != nil {
		assert.ErrorIs(t, err, OutOfStock)
	}

	assert.Equal(t, 90, inStockBatch.AvailableQuantity())
	assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
}

func TestPrefersEarlierBatches(t *testing.T) {
	earliest := NewBatch("speedy-batch", "MINIMALIST-SPOON", 100, now())
	medium := NewBatch("normal-batch", "MINIMALIST-SPOON", 100, tomorrow())
	latest := NewBatch("slow-batch", "MINIMALIST-SPOON", 100, later())

	line := NewOrderLine("order1", "MINIMALIST-SPOON", 10)

	_, err := Allocate(line, []*Batch{medium, earliest, latest})
	if err != nil {
		assert.ErrorIs(t, err, OutOfStock)
	}

	assert.Equal(t, 90, earliest.AvailableQuantity())
	assert.Equal(t, 100, medium.AvailableQuantity())
	assert.Equal(t, 100, latest.AvailableQuantity())
}

func TestReturnsAllocatedBatchRef(t *testing.T) {
	inStockBatch := NewBatch("in-stock-batch-ref", "HIGHBROW-POSTER", 100, time.Time{})
	shipmentBatch := NewBatch("shipment-batch-ref", "HIGHBROW-POSTER", 100, tomorrow())
	line := NewOrderLine("oref", "HIGHBROW-POSTER", 10)

	ref, err := Allocate(line, []*Batch{inStockBatch, shipmentBatch})
	if err != nil {
		assert.ErrorIs(t, err, OutOfStock)
	}

	assert.Equal(t, inStockBatch.Reference, ref.Reference)
}

//func TestReturnsOutOfStockExceptionIfCannotAllocate(t *testing.T) {
//	batch := NewBatch("batch1", "SMALL-FORK", 10, time.Time{})
//	line1 := NewOrderLine("order1", "SMALL-FORK", 10)
//	line2 := NewOrderLine("order2", "SMALL-FORK", 10)
//
//	_, err := Allocate(line1, []*Batch{batch})
//	if err != nil {
//		assert.ErrorIs(t, err, OutOfStock)
//	}
//
//	_, err = Allocate(line2, []*Batch{batch})
//	assert.ErrorIs(t, err, OutOfStock)
//}

func TestReturnsOutOfStockEventIfCannotAllocate(t *testing.T) {
	sku := "SMALL-FORK"
	batch := NewBatch("batch1", sku, 10, time.Time{})
	product := NewProduct(sku, []*Batch{batch})
	line1 := NewOrderLine("order1", sku, 10)
	line2 := NewOrderLine("order2", sku, 10)

	_ = product.Allocate(line1)
	allocation := product.Allocate(line2)
	assert.True(t, product.HasOutOfStockEventAsLast())
	assert.NotNil(t, allocation)
}
