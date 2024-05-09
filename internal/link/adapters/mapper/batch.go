package mapper

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func BatchToDomain(bch *entity.Batch) *domain.Batch {
	batch := domain.NewBatch(bch.Reference, bch.SKU, bch.PurchasedQuantity, bch.ETA)
	batch.ID = bch.ID
	for _, line := range bch.Allocations {
		batch.Allocate(OrderLineToDomain(line.OrderLine))
	}
	return batch
}

func BatchToEntity(bch *domain.Batch) *entity.Batch {
	batch := &entity.Batch{
		ID:                bch.ID,
		Reference:         bch.Reference,
		SKU:               bch.SKU,
		PurchasedQuantity: bch.PurchasedQuantity,
		ETA:               bch.ETA,
	}
	for _, orderLine := range bch.Allocations() {
		line := OrderLineToEntity(orderLine)
		batch.Allocations = append(batch.Allocations, ConvertAllocation(batch, line))
	}
	return batch
}

func BatchToDomainMany(bchList []*entity.Batch) []*domain.Batch {
	batches := make([]*domain.Batch, len(bchList), len(bchList))
	for i, line := range bchList {
		batches[i] = BatchToDomain(line)
	}
	return batches
}

func BatchToEntityMany(bchList []*domain.Batch) []*entity.Batch {
	batches := make([]*entity.Batch, len(bchList), len(bchList))
	for i, line := range bchList {
		batches[i] = BatchToEntity(line)
	}
	return batches
}

func BatchToArrayOfPointers(list []entity.Batch) []*entity.Batch {
	batches := make([]*entity.Batch, len(list), len(list))
	for i := range list {
		batches[i] = &list[i]
	}
	return batches
}

func BatchToArrayOfValues(list []*entity.Batch) []entity.Batch {
	batches := make([]entity.Batch, len(list), len(list))
	for i := range list {
		batches[i] = *list[i]
	}
	return batches
}
