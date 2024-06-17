package mapper

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/model"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func OrderLineToDomain(line entity.OrderLine) model.OrderLine {
	return model.OrderLine{
		ID:      line.ID,
		OrderID: line.OrderID,
		SKU:     line.SKU,
		Qty:     line.Qty,
	}
}

func OrderLineToEntity(line model.OrderLine) entity.OrderLine {
	return entity.OrderLine{
		ID:      line.ID,
		OrderID: line.OrderID,
		SKU:     line.SKU,
		Qty:     line.Qty,
	}
}

func OrderLineToDomainMany(lines []entity.OrderLine) []model.OrderLine {
	domainLines := make([]model.OrderLine, len(lines), len(lines))
	for i, line := range lines {
		domainLines[i] = OrderLineToDomain(line)
	}
	return domainLines
}

func OrderLineToEntityMany(lines []model.OrderLine) []entity.OrderLine {
	domainLines := make([]entity.OrderLine, len(lines), len(lines))
	for i, line := range lines {
		domainLines[i] = OrderLineToEntity(line)
	}
	return domainLines
}
