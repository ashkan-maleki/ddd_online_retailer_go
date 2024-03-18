package mapper

import "github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"

func ConvertAllocation(batch *entity.Batch, line entity.OrderLine) entity.Allocation {
	return entity.Allocation{
		OrderLineID: line.ID,
		OrderLine:   line,
		BatchID:     batch.ID,
	}
}
