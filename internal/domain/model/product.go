package model

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain_events"
	"sort"
)

type Product struct {
	*domain.BaseEntity
	SKU           string
	VersionNumber int
	Batches       []*Batch
}

func (p *Product) HasOutOfStockEventAsLast() bool {
	last := p.LastEvent()
	ev := last.(*domain_events.OutOfStock)
	return ev.Sku() == p.SKU
}

func NewProduct(SKU string, batches []*Batch) *Product {
	return &Product{SKU: SKU, Batches: batches, VersionNumber: 0, BaseEntity: domain.NewBaseEntity()}
}

func (p *Product) Allocate(line OrderLine) (*Batch, error) {
	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].ETA.Before(p.Batches[j].ETA)
	})
	for _, batch := range p.Batches {

		if batch.CanAllocate(line) {
			batch.Allocate(line)
			p.VersionNumber += 1
			p.AddDomainEvent(domain_events.NewAllocated(line.OrderID, line.SKU, batch.Reference, line.Qty))
			return batch, nil
		}
	}

	p.AddDomainEvent(domain_events.NewOutOfStockEvent(line.SKU))
	return nil, OutOfStockErr
}

func (p *Product) ChangeBatchQuantity(ref string, qty int) {
	var batch *Batch
	for _, b := range p.Batches {
		if b.Reference == ref {
			batch = b
		}
	}
	batch.PurchasedQuantity = qty

	for batch.AvailableQuantity() < 0 {
		line := batch.DeallocateOne()
		p.AddDomainEvent(domain_events.NewDeallocated(line.OrderID, line.SKU, line.Qty))
	}
}
