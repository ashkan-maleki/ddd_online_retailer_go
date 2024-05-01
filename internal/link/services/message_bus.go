package services

import (
	"context"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters/mapper"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

type HandleFunc func(context.Context, events.Event, *adapters.ProductRepo) (any, error)

var Handlers = make(map[string][]HandleFunc)

func SendOutOfStockNotification(_ context.Context, event events.Event, _ *adapters.ProductRepo) (any, error) {
	outOfStock, ok := event.(events.OutOfStock)
	if ok {
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}
	SendEmail("stock@eshop.com", fmt.Sprintf("out of stock for %v", outOfStock.Sku()))
	return nil, nil
}

func AddBatch(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	batchCreated, ok := event.(events.BatchCreated)
	if ok {
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	product := repo.Get(ctx, batchCreated.Sku())
	if product == nil {
		product = &entity.Product{
			SKU:     batchCreated.Sku(),
			Batches: make([]entity.Batch, 0),
		}
	}

	batch := entity.Batch{
		Reference:         batchCreated.Ref(),
		SKU:               batchCreated.Sku(),
		PurchasedQuantity: batchCreated.Qty(),
		ETA:               batchCreated.Eta(),
	}
	product.Batches = append(product.Batches, batch)
	err := repo.Add(ctx, product)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func Allocate(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	allocationRequired, ok := event.(events.AllocationRequired)
	if ok {
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	sku := allocationRequired.Sku()
	line := domain.NewOrderLine(allocationRequired.OrderId(), sku, allocationRequired.Qty())

	product := mapper.ProductToDomain(repo.Get(ctx, sku))
	if product == nil {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}

	batch, err := product.Allocate(line)
	if err != nil {
		return "", err
	}

	err = repo.Update(ctx, mapper.ProductToEntity(product))
	if err != nil {
		return nil, err
	}
	return batch.Reference, nil
}

func ChangeBatchQuantity(ctx context.Context, event events.Event, repo *adapters.ProductRepo) (any, error) {
	batchQuantityChanged, ok := event.(events.BatchQuantityChanged)
	if ok {
		return nil, fmt.Errorf("wrong event type %v", event.Name())
	}

	fmt.Println(batchQuantityChanged)

	return nil, nil
}

func Register() {
	Handlers[events.BatchCreatedEvent] = []HandleFunc{AddBatch}
	Handlers[events.BatchQuantityChangedEvent] = []HandleFunc{ChangeBatchQuantity}
	Handlers[events.AllocationRequiredEvent] = []HandleFunc{Allocate}
	Handlers[events.OutOfStockEvent] = []HandleFunc{SendOutOfStockNotification}
}

func Handle(ctx context.Context, event events.Event) ([]any, error) {
	handlers, ok := Handlers[event.Name()]
	if !ok {
		return nil, fmt.Errorf("no handler is registered for %v", event.Name())
	}
	repo, err := adapters.NewProductRepo()
	if err != nil {
		return nil, err
	}
	results := make([]any, 0)
	for _, handler := range handlers {
		result, err := handler(ctx, event, repo)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
