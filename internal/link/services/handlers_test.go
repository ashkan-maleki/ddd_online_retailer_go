package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters/mapper"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func newBatchRepo() *adapters.ProductRepo {
	repo, err := adapters.NewProductRepo()
	if err != nil {
		log.Println(err)
		panic("new batch repo failed")
	}
	return repo
}

func TestAddBatch_NewProduct(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "b1"
	sku := "CRUNCHY-ARMCHAIR"

	_, err := Handle(ctx, events.NewBatchCreated(sku, ref, 100, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	product := uow.Product.Get(ctx, sku)
	assert.NotNil(t, product)
	assert.NotNil(t, product.Batches)
	found := false
	for _, batch := range product.Batches {
		if batch.Reference == ref {
			found = true
		}
	}
	assert.True(t, found)
}

func TestAddBatch_ExistingProduct(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "b2"
	sku := "GARISH-RUG"

	_, err := Handle(ctx, events.NewBatchCreated(sku, "b1", 100, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "1st handle function error: "+err.Error())
	}

	_, err = Handle(ctx, events.NewBatchCreated(sku, ref, 99, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "2nd handle function error: "+err.Error())
	}

	product := uow.Product.Get(ctx, sku)
	assert.NotNil(t, product)
	assert.NotNil(t, product.Batches)
	found := false
	for _, batch := range product.Batches {
		if batch.Reference == ref {
			found = true
		}
	}
	assert.True(t, found)
}

func TestAllocate_ReturnsAllocation(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "batch1"
	sku := "COMPLICATED-LAMP"

	_, err := Handle(ctx, events.NewBatchCreated(sku, ref, 100, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	results, err := Handle(ctx, events.NewAllocationRequired("o1", sku, 10), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}
	assert.Equal(t, ref, results[0])
}

func TestAllocate_ReturnsInvalidSku(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "b1"
	sku := "AREALSKU"

	_, err := Handle(ctx, events.NewBatchCreated(sku, ref, 100, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	_, err = Handle(ctx, events.NewAllocationRequired("o1", "NONEXISTENTSKU", 10), uow.Product)
	if err != nil {
		assert.ErrorIs(t, err, InvalidSku)
	}
}

func TestAllocate_SendsEmailOnOutOfStockError(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "b1"
	sku := "POPULAR-CURTAINS"

	_, err := Handle(ctx, events.NewBatchCreated(sku, ref, 9, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	results, err := Handle(ctx, events.NewAllocationRequired("o1", sku, 10), uow.Product)
	if err != nil && !errors.Is(err, domain.OutOfStockErr) {
		assert.Fail(t, "handle function error: "+err.Error())
	}
	fmt.Println("tests: ", len(results))
	fmt.Println("tests res 0: ", results[0])
	assert.Equal(t, fmt.Sprintf("out of stock for %v", sku), results[0])
}

func TestChangeBatchQuantity_ChangesAvailableQuantity(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref := "batch1"
	sku := "ADORABLE-SETTEE"

	_, err := Handle(ctx, events.NewBatchCreated(sku, ref, 100, time.Time{}), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	product := uow.Product.Get(ctx, sku)
	assert.NotNil(t, product)
	assert.Len(t, product.Batches, 1)
	domainProduct := mapper.ProductToDomain(product)
	fmt.Println("step 1")
	assert.Equal(t, 100, domainProduct.Batches[0].AvailableQuantity())

	_, err = Handle(ctx, events.NewBatchQuantityChanged(ref, 50), uow.Product)
	if err != nil {
		assert.Fail(t, "handle function error: "+err.Error())
	}

	product1 := uow.Product.Get(ctx, sku)
	assert.NotNil(t, product1)
	assert.Len(t, product1.Batches, 1)
	domainProduct = mapper.ProductToDomain(product1)
	fmt.Println("step 2")
	assert.Equal(t, 50, domainProduct.Batches[0].AvailableQuantity())
}

func TestChangeBatchQuantity_ReallocatesIfNecessary(t *testing.T) {
	ctx := context.Background()
	uow := UoW()

	ref1 := "batch1"
	ref2 := "batch2"
	sku := "INDIFFERENT-TABLE"

	eventsList := []events.Event{
		events.NewBatchCreated(sku, ref1, 50, time.Time{}),
		events.NewBatchCreated(sku, ref2, 50, time.Now()),
		events.NewAllocationRequired("order1", sku, 20),
		events.NewAllocationRequired("order2", sku, 20),
	}

	for _, event := range eventsList {
		_, err := Handle(ctx, event, uow.Product)
		if err != nil {
			assert.Fail(t, "handle function error: "+err.Error())
		}
	}

	batches := uow.Product.Get(ctx, sku).Batches
	for i, batch := range batches {
		fmt.Printf("#%d: number of allocated (ent) is %v\n", i, len(batch.Allocations))
		toDomain := mapper.BatchToDomain(&batch)
		fmt.Printf("#%d: quantity is %v\n", i, toDomain.AvailableQuantity())
		fmt.Printf("#%d: number of allocated (domain) is %v\n", i, len(toDomain.Allocations()))
	}
	assert.Len(t, batches, 2)
	assert.Equal(t, 10, mapper.BatchToDomain(&batches[1]).AvailableQuantity())
	assert.Equal(t, 50, mapper.BatchToDomain(&batches[0]).AvailableQuantity())

	//_, err := Handle(ctx, events.NewBatchQuantityChanged(ref1, 25), uow.Product)
	//if err != nil {
	//	assert.Fail(t, "handle function error: "+err.Error())
	//}
	//
	//batches = uow.Product.Get(ctx, sku).Batches
	//assert.Equal(t, 5, mapper.BatchToDomain(&batches[0]).AvailableQuantity())
	//assert.Equal(t, 30, mapper.BatchToDomain(&batches[1]).AvailableQuantity())
}
