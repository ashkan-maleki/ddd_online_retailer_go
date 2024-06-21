package model

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain_events"
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

	product.Allocate(line)
	assert.Len(t, product.DomainEvents(), 1)
	AssertAllocatedEvent(t, product.DomainEvents()[0], line, inStockBatch)

	assert.Equal(t, 90, inStockBatch.AvailableQuantity())
	assert.Equal(t, 100, shipmentBatch.AvailableQuantity())
}

func AssertAllocatedEvent(t *testing.T, eve domain.Event, line OrderLine, inStockBatch *Batch) {
	event, err := domain_events.ConvertAllocate(eve)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, event.OrderID(), line.OrderID)
	assert.Equal(t, event.SKU(), line.SKU)
	assert.Equal(t, event.Qty(), line.Qty)
	assert.Equal(t, event.BatchRef(), inStockBatch.Reference)
}

func AssertOutOfStockEvent(t *testing.T, eve domain.Event, sku string) {
	event, err := domain_events.ConvertOutOfStock(eve)
	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, event.SKU(), sku)
}

func TestPrefersEarlierBatches(t *testing.T) {
	sku := "MINIMALIST-SPOON"
	earliest := NewBatch("speedy-batch", sku, 100, now())
	medium := NewBatch("normal-batch", sku, 100, tomorrow())
	latest := NewBatch("slow-batch", sku, 100, later())
	product := NewProduct(sku, []*Batch{medium, earliest, latest})
	line := NewOrderLine("order1", sku, 10)

	product.Allocate(line)
	assert.Len(t, product.DomainEvents(), 1)
	AssertAllocatedEvent(t, product.DomainEvents()[0], line, earliest)

	assert.Equal(t, 90, earliest.AvailableQuantity())
	assert.Equal(t, 100, medium.AvailableQuantity())
	assert.Equal(t, 100, latest.AvailableQuantity())
}

//func TestReturnsAllocatedBatchRef(t *testing.T) {
//	sku := "HIGHBROW-POSTER"
//	inStockBatch := NewBatch("in-stock-batch-ref", sku, 100, time.Time{})
//	shipmentBatch := NewBatch("shipment-batch-ref", sku, 100, tomorrow())
//	product := NewProduct(sku, []*Batch{inStockBatch, shipmentBatch})
//	line := NewOrderLine("oref", sku, 10)
//
//	ref, err := product.AllocateDeprecated(line)
//	if err != nil {
//		assert.ErrorIs(t, err, OutOfStockErr)
//	}
//
//	assert.Equal(t, inStockBatch.Reference, ref.Reference)
//}

// //func TestReturnsOutOfStockExceptionIfCannotAllocate(t *testing.T) {
// //	batch := NewBatch("batch1", "SMALL-FORK", 10, time.Time{})
// //	line1 := NewOrderLine("order1", "SMALL-FORK", 10)
// //	line2 := NewOrderLine("order2", "SMALL-FORK", 10)
// //
// //	_, err := AllocateDeprecated(line1, []*Batch{batch})
// //	if err != nil {
// //		assert.ErrorIs(t, err, OutOfStockErr)
// //	}
// //
// //	_, err = AllocateDeprecated(line2, []*Batch{batch})
// //	assert.ErrorIs(t, err, OutOfStockErr)
// //}
func TestRecordsOutOfStockEventIfCannotAllocate(t *testing.T) {
	sku := "SMALL-FORK"
	batch := NewBatch("batch1", sku, 10, time.Time{})
	product := NewProduct(sku, []*Batch{batch})
	line1 := NewOrderLine("order1", sku, 10)
	line2 := NewOrderLine("order2", sku, 10)

	product.Allocate(line1)
	product.Allocate(line2)
	assert.Len(t, product.DomainEvents(), 2)
	AssertAllocatedEvent(t, product.DomainEvents()[0], line1, batch)
	AssertOutOfStockEvent(t, product.DomainEvents()[1], sku)
}
