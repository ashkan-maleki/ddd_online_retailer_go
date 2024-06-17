package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/model"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters/mapper"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
	"time"
)

type ProductService struct {
	repo *adapters.ProductRepo
}

func NewBatchService(repo *adapters.ProductRepo) *ProductService {
	return &ProductService{repo: repo}
}

var InvalidSku = errors.New("invalid sku")

func IsValidSku(sku string, batches []*entity.Batch) bool {
	for _, batch := range batches {
		if sku == batch.SKU {
			return true
		}
	}
	return false
}

func (service *ProductService) AddBatch(ctx context.Context, reference, sku string, qty int, eta time.Time) error {
	batch := &model.Batch{
		Reference:         reference,
		SKU:               sku,
		PurchasedQuantity: qty,
		ETA:               eta,
	}

	product := mapper.ProductToDomain(service.repo.Get(ctx, sku))
	if product == nil {
		product = model.NewProduct(sku, make([]*model.Batch, 0))
	}
	product.Batches = append(product.Batches, batch)
	productEntity := mapper.ProductToEntity(product)
	return service.repo.Add(ctx, productEntity)
}

func (service *ProductService) Allocate(ctx context.Context, orderID, sku string, qty int) (string, error) {
	line := model.NewOrderLine(orderID, sku, qty)

	product := mapper.ProductToDomain(service.repo.Get(ctx, sku))
	if product == nil {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}

	batch, err := product.Allocate(line)
	if err != nil {
		return "", err
	}
	// TODO: Save allocations in the database
	return batch.Reference, nil
}
