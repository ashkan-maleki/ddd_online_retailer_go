package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
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
