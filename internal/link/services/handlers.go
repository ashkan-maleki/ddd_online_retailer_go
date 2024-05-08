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

	product := mapper.ProductToDomain(repo.Get(ctx, sku))
	if product == nil {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}
	batch, allocationErr := product.Allocate(line)

	if allocationErr == nil || errors.Is(allocationErr, domain.OutOfStockErr) {
		productEntity := mapper.ProductToEntity(product)
		//fmt.Println("productEntity: ", productEntity)
		//fmt.Println("productEntity sku: ", productEntity.SKU)
		//fmt.Println("productEntity events size: ", len(productEntity.Events()))
		//if len(productEntity.Events()) > 0 {
		//	fmt.Println("productEntity event: ", productEntity.Events()[0])
		//}
		err := repo.Update(ctx, productEntity)
		if err != nil {
			return nil, err
		}
		//fmt.Println("repo: ", collectNewEvents(repo))
	}

	if allocationErr != nil {
		return nil, allocationErr
	}
	return batch.Reference, nil
}

func ChangeBatchQuantity(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	batchQuantityChanged, ok := event.(events.BatchQuantityChanged)
	if ok {
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	productEnt := repo.GetByBatchRef(ctx, batchQuantityChanged.Ref())
	product := mapper.ProductToDomain(productEnt)
	product.ChangeBatchQuantity(batchQuantityChanged.Ref(), batchQuantityChanged.Qty())
	err := repo.Update(ctx, mapper.ProductToEntity(product))
	if err != nil {
		return nil, err
	}
	return nil, nil
}
