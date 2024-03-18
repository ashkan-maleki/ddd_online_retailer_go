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

type BatchService struct {
	repo *adapters.BatchRepo
}

func NewBatchService(repo *adapters.BatchRepo) *BatchService {
	return &BatchService{repo: repo}
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

func (service *BatchService) AddBatch(ctx context.Context, reference, sku string, qty int, eta time.Time) {
	service.repo.Add(ctx, &entity.Batch{
		Reference:         reference,
		SKU:               sku,
		PurchasedQuantity: qty,
		ETA:               eta,
	})
}

func (service *BatchService) Allocate(ctx context.Context, orderID, sku string, qty int) (string, error) {
	line := domain.NewOrderLine(orderID, sku, qty)
	batches := service.repo.List(ctx)
	if !IsValidSku(sku, batches) {
		return "", fmt.Errorf("sku validation: %w", InvalidSku)
	}
	batch, err := domain.Allocate(line, mapper.BatchToDomainMany(batches))
	if err != nil {
		return "", fmt.Errorf("domain allocation: %w", err)
	}

	return batch.Reference, nil
}
