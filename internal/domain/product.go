package domain

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain_events"
	"sort"
)

type BaseEntity struct {
	events []domain_events.Event
}

func NewBaseEntity() *BaseEntity {
	return &BaseEntity{events: make([]domain_events.Event, 0)}
}

func (p *BaseEntity) Events() []domain_events.Event {
	return p.events
}

func (p *BaseEntity) AddEvent(event domain_events.Event) {
	p.events = append(p.events, event)
}

func (p *BaseEntity) HasEvent() bool {
	return len(p.events) > 0
}

type Product struct {
	*BaseEntity
	SKU           string
	VersionNumber int
	Batches       []*Batch
}

func (p *Product) PopEvent() domain_events.Event {
	if len(p.events) > 1 {
		p.events = p.events[1:]
		return p.events[0]
	} else if len(p.events) == 1 {
		p.events = make([]domain_events.Event, 0)
		return p.events[0]
	} else {
		return nil
	}
}

func (p *Product) HasOutOfStockEventAsLast() bool {
	if len(p.events) == 0 {
		return false
	}
	last := p.events[len(p.events)-1]
	ev := last.(*domain_events.OutOfStock)
	return ev.Sku() == p.SKU
}

func NewProduct(SKU string, batches []*Batch) *Product {
	return &Product{SKU: SKU, Batches: batches, VersionNumber: 0, BaseEntity: NewBaseEntity()}
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

	p.events = append(p.events, domain_events.NewOutOfStockEvent(line.SKU))
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
		p.events = append(p.events, domain_events.NewAllocationRequired(line.OrderID, line.SKU, line.Qty))
	}
}
