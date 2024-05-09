package mapper

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func OrderLineToDomain(line entity.OrderLine) domain.OrderLine {
	return domain.OrderLine{
		ID:      line.ID,
		OrderID: line.OrderID,
		SKU:     line.SKU,
		Qty:     line.Qty,
	}
}

func OrderLineToEntity(line domain.OrderLine) entity.OrderLine {
	return entity.OrderLine{
		ID:      line.ID,
		OrderID: line.OrderID,
		SKU:     line.SKU,
		Qty:     line.Qty,
	}
}

func OrderLineToDomainMany(lines []entity.OrderLine) []domain.OrderLine {
	domainLines := make([]domain.OrderLine, len(lines), len(lines))
	for i, line := range lines {
		domainLines[i] = OrderLineToDomain(line)
	}
	return domainLines
}

func OrderLineToEntityMany(lines []domain.OrderLine) []entity.OrderLine {
	domainLines := make([]entity.OrderLine, len(lines), len(lines))
	for i, line := range lines {
		domainLines[i] = OrderLineToEntity(line)
	}
	return domainLines
}
