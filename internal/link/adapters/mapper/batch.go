package mapper

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/model"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func BatchToDomain(bch *entity.Batch) *model.Batch {
	batch := model.NewBatch(bch.Reference, bch.SKU, bch.PurchasedQuantity, bch.ETA)
	batch.ID = bch.ID
	fmt.Printf(">>> (mapper), allocation lenth for entity is %d\n", len(bch.Allocations))
	for _, line := range bch.Allocations {
		fmt.Printf("loop startd: ")
		fmt.Printf("%v\n", line)
		batch.Allocate(OrderLineToDomain(line.OrderLine))
		fmt.Printf("loop ended \n")
	}
	fmt.Printf(">>> (mapper), allocation lenth for domain is %d\n", len(batch.Allocations()))
	return batch
}

func BatchToEntity(bch *model.Batch) *entity.Batch {
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

func BatchToDomainMany(bchList []*entity.Batch) []*model.Batch {
	batches := make([]*model.Batch, len(bchList), len(bchList))
	for i, line := range bchList {
		batches[i] = BatchToDomain(line)
	}
	return batches
}

func BatchToEntityMany(bchList []*model.Batch) []*entity.Batch {
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
