package mapper

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/persistence/entity"
)

func ProductToDomain(p *entity.Product) *domain.Product {
	batches := BatchToDomainMany(BatchToArrayOfPointers(p.Batches))
	product := domain.NewProduct(p.SKU, batches)
	product.VersionNumber = p.VersionNumber
	return product
}

func ProductToEntity(p *domain.Product) *entity.Product {
	batches := BatchToArrayOfValues(BatchToEntityMany(p.Batches))
	product := &entity.Product{
		SKU:           p.SKU,
		Batches:       batches,
		VersionNumber: p.VersionNumber,
	}
	return product
}

func ProductToDomainMany(list []*entity.Product) []*domain.Product {
	products := make([]*domain.Product, len(list), len(list))
	for i, line := range list {
		products[i] = ProductToDomain(line)
	}
	return products
}

func ProductToEntityMany(list []*domain.Product) []*entity.Product {
	products := make([]*entity.Product, len(list), len(list))
	for i, line := range list {
		products[i] = ProductToEntity(line)
	}
	return products
}
