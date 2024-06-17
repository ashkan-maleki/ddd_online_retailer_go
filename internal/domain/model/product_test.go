package model

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

func TestPrefersWarehouseBatchesToShipments(t *testing.T) {
	sku := "RETRO-CLOCK"
	inStockBatch := NewBatch("in-stock-batch", sku, 100, time.Time{})
	shipmentBatch := NewBatch("shipment-batch", sku, 100, tomorrow())
	product := NewProduct(sku, []*Batch{inStockBatch, shipmentBatch})
	line := NewOrderLine("oref", sku, 10)

	_, err := product.Allocate(line)
	if err != nil {
		assert.ErrorIs(t, err, OutOfStockErr)
	}

	assert.Equal(t, 90, inStockBatch.AvailableQuantity())
	assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
}

func TestPrefersEarlierBatches(t *testing.T) {
	sku := "MINIMALIST-SPOON"
	earliest := NewBatch("speedy-batch", sku, 100, now())
	medium := NewBatch("normal-batch", sku, 100, tomorrow())
	latest := NewBatch("slow-batch", sku, 100, later())
	product := NewProduct(sku, []*Batch{medium, earliest, latest})
	line := NewOrderLine("order1", sku, 10)

	_, err := product.Allocate(line)
	if err != nil {
		assert.ErrorIs(t, err, OutOfStockErr)
	}

	assert.Equal(t, 90, earliest.AvailableQuantity())
	assert.Equal(t, 100, medium.AvailableQuantity())
	assert.Equal(t, 100, latest.AvailableQuantity())
}

func TestReturnsAllocatedBatchRef(t *testing.T) {
	sku := "HIGHBROW-POSTER"
	inStockBatch := NewBatch("in-stock-batch-ref", sku, 100, time.Time{})
	shipmentBatch := NewBatch("shipment-batch-ref", sku, 100, tomorrow())
	product := NewProduct(sku, []*Batch{inStockBatch, shipmentBatch})
	line := NewOrderLine("oref", sku, 10)

	ref, err := product.Allocate(line)
	if err != nil {
		assert.ErrorIs(t, err, OutOfStockErr)
	}

	assert.Equal(t, inStockBatch.Reference, ref.Reference)
}

//
////func TestReturnsOutOfStockExceptionIfCannotAllocate(t *testing.T) {
////	batch := NewBatch("batch1", "SMALL-FORK", 10, time.Time{})
////	line1 := NewOrderLine("order1", "SMALL-FORK", 10)
////	line2 := NewOrderLine("order2", "SMALL-FORK", 10)
////
////	_, err := Allocate(line1, []*Batch{batch})
////	if err != nil {
////		assert.ErrorIs(t, err, OutOfStockErr)
////	}
////
////	_, err = Allocate(line2, []*Batch{batch})
////	assert.ErrorIs(t, err, OutOfStockErr)
////}
//
//func TestReturnsOutOfStockEventIfCannotAllocate(t *testing.T) {
//	sku := "SMALL-FORK"
//	batch := NewBatch("batch1", sku, 10, time.Time{})
//	product := NewProduct(sku, []*Batch{batch})
//	line1 := NewOrderLine("order1", sku, 10)
//	line2 := NewOrderLine("order2", sku, 10)
//
//	_, _ = product.Allocate(line1)
//	allocation, err := product.Allocate(line2)
//	assert.True(t, product.HasOutOfStockEventAsLast())
//	assert.Nil(t, allocation)
//	assert.NotNil(t, err)
//	assert.ErrorIs(t, err, OutOfStockErr)
//}
