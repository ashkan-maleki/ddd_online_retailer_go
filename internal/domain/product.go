package domain

import "sort"

type Product struct {
	SKU           string
	VersionNumber int
	Batches       []*Batch
}

func NewProduct(SKU string, batches []*Batch) *Product {
	return &Product{SKU: SKU, Batches: batches, VersionNumber: 0}
}

func (p *Product) Allocate(line OrderLine) (*Batch, error) {
	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].ETA.Before(p.Batches[j].ETA)
	})
	for _, batch := range p.Batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			return batch, nil
		}
	}
	return nil, OutOfStock
}
