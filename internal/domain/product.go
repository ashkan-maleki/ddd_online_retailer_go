package domain

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain_events"
	"sort"
)

type Entity interface {
	DomainEvents() []domain_events.Event
	AddDomainEvent(domain_events.Event)
	HasDomainEvent() bool
}

type BaseEntity struct {
	domainEvents []domain_events.Event
}

var _ Entity = (*BaseEntity)(nil)

func NewBaseEntity() *BaseEntity {
	return &BaseEntity{domainEvents: make([]domain_events.Event, 0)}
}

func (p *BaseEntity) DomainEvents() []domain_events.Event {
	return p.domainEvents
}

func (p *BaseEntity) AddDomainEvent(event domain_events.Event) {
	p.domainEvents = append(p.domainEvents, event)
}

func (p *BaseEntity) HasDomainEvent() bool {
	return len(p.domainEvents) > 0
}

func (p *BaseEntity) PopEvent() domain_events.Event {
	if len(p.domainEvents) > 1 {
		p.domainEvents = p.domainEvents[1:]
		return p.domainEvents[0]
	} else if len(p.domainEvents) == 1 {
		p.domainEvents = make([]domain_events.Event, 0)
		return p.domainEvents[0]
	} else {
		return nil
	}
}

type Product struct {
	*BaseEntity
	SKU           string
	VersionNumber int
	Batches       []*Batch
}

func (p *Product) HasOutOfStockEventAsLast() bool {
	if len(p.domainEvents) == 0 {
		return false
	}
	last := p.domainEvents[len(p.domainEvents)-1]
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
		p.AddDomainEvent(domain_events.NewAllocationRequired(line.OrderID, line.SKU, line.Qty))
	}
}
