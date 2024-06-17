package mapper

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/model"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func ProductToDomain(p *entity.Product) *model.Product {
	if p == nil {
		return nil
	}
	batches := BatchToDomainMany(BatchToArrayOfPointers(p.Batches))
	product := model.NewProduct(p.SKU, batches)
	for _, event := range p.Events() {
		product.AddDomainEvent(event)
	}
	product.VersionNumber = p.VersionNumber
	return product
}

func ProductToEntity(p *model.Product) *entity.Product {
	batches := BatchToArrayOfValues(BatchToEntityMany(p.Batches))
	product := &entity.Product{
		SKU:           p.SKU,
		Batches:       batches,
		VersionNumber: p.VersionNumber,
	}
	for _, event := range p.DomainEvents() {
		product.AddEvent(event)
	}
	return product
}

func ProductToDomainMany(list []*entity.Product) []*model.Product {
	products := make([]*model.Product, len(list), len(list))
	for i, line := range list {
		products[i] = ProductToDomain(line)
	}
	return products
}

func ProductToEntityMany(list []*model.Product) []*entity.Product {
	products := make([]*entity.Product, len(list), len(list))
	for i, line := range list {
		products[i] = ProductToEntity(line)
	}
	return products
}
