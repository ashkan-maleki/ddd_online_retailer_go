package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters/mapper"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func SendOutOfStockNotification(_ context.Context, event events.Event, _ *adapters.ProductRepo) (any, error) {

	var outOfStock *events.OutOfStock
	switch a := event.(type) {
	case *events.OutOfStock:
		outOfStock = a
		break

	default:
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	emailMessage := fmt.Sprintf("out of stock for %v", outOfStock.Sku())
	SendEmail("stock@eshop.com", emailMessage)
	fmt.Println("email is  sent")
	return emailMessage, nil
}

func AddBatch(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	var batchCreated *events.BatchCreated
	switch a := event.(type) {
	case *events.BatchCreated:
		batchCreated = a
		break

	default:
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	batch := entity.Batch{
		Reference:         batchCreated.Ref(),
		SKU:               batchCreated.GetSku(),
		PurchasedQuantity: batchCreated.Qty(),
		ETA:               batchCreated.Eta(),
	}

	product := repo.Get(ctx, batchCreated.GetSku())
	if product == nil {
		product = &entity.Product{
			SKU:     batchCreated.GetSku(),
			Batches: make([]entity.Batch, 0),
		}
		product.Batches = append(product.Batches, batch)
		err := repo.Add(ctx, product)
		if err != nil {
			return nil, err
		}
	} else {
		product.Batches = append(product.Batches, batch)
		err := repo.Update(ctx, product)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func Allocate(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	var allocationRequired *events.AllocationRequired
	switch a := event.(type) {
	case *events.AllocationRequired:
		allocationRequired = a
		break
	default:
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	sku := allocationRequired.Sku()
	line := domain.NewOrderLine(allocationRequired.OrderId(), sku, allocationRequired.Qty())

	get := repo.Get(ctx, sku)
	gotBatch := get.Batches[0]
	fmt.Printf("Batch with ref %v allocated size is %d\n", gotBatch.Reference, len(gotBatch.Allocations))
	product := mapper.ProductToDomain(get)
	fmt.Printf("Batch with ref %v allocated size is %d\n", product.Batches[1].Reference, len(product.Batches[1].Allocations()))
	if product == nil {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}
	batch, allocationErr := product.Allocate(line)
	fmt.Printf("OrderID: %v, Product with ref %v has %d available\n", line.OrderID, batch.Reference, batch.AvailableQuantity())

	if allocationErr == nil || errors.Is(allocationErr, domain.OutOfStockErr) {
		productEntity := mapper.ProductToEntity(product)
		err := repo.Update(ctx, productEntity)
		if err != nil {
			return nil, err
		}
	}

	if allocationErr != nil {
		return nil, allocationErr
	}
	return batch.Reference, nil
}

func ChangeBatchQuantity(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	var batchQuantityChanged *events.BatchQuantityChanged
	switch a := event.(type) {
	case *events.BatchQuantityChanged:
		batchQuantityChanged = a
		break

	default:
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	productEnt := repo.GetByBatchRef(ctx, batchQuantityChanged.Ref())

	product := mapper.ProductToDomain(productEnt)
	product.ChangeBatchQuantity(batchQuantityChanged.Ref(), batchQuantityChanged.Qty())

	toEntity := mapper.ProductToEntity(product)

	err := repo.Update(ctx, toEntity)

	if err != nil {
		return nil, err
	}
	return nil, nil
}
