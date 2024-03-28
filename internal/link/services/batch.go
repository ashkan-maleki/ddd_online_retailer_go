package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
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

func (service *ProductService) AddBatch(ctx context.Context, reference, sku string, qty int, eta time.Time) {
	batch := &domain.Batch{
		Reference:         reference,
		SKU:               sku,
		PurchasedQuantity: qty,
		ETA:               eta,
	}
	product := mapper.ProductToDomain(service.repo.Get(ctx, sku))
	if product == nil {
		product = domain.NewProduct(sku, make([]*domain.Batch, 0))
	}
	product.Batches = append(product.Batches, batch)
	service.repo.Add(ctx, mapper.ProductToEntity(product))

}

func (service *ProductService) Allocate(ctx context.Context, orderID, sku string, qty int) (string, error) {
	line := domain.NewOrderLine(orderID, sku, qty)
	batches := mapper.BatchToArrayOfPointers(service.repo.Get(ctx, sku).Batches)
	if !IsValidSku(sku, batches) {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}
	batch, err := domain.Allocate(line, mapper.BatchToDomainMany(batches))
	if err != nil {
		return "", fmt.Errorf("domain allocation: %w", err)
	}
	return batch.Reference, nil
}
