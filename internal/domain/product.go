package domain

import (
	"sort"
)

type Product struct {
	SKU           string
	VersionNumber int
	Batches       []*Batch
	events        []Event
}

func (p *Product) Events() []Event {
	return p.events
}

func (p *Product) AddEvent(event Event) {
	p.events = append(p.events, event)
}

func (p *Product) HasEvent() bool {
	return len(p.events) > 0
}

func (p *Product) PopEvent() Event {
	if len(p.events) > 1 {
		p.events = p.events[1:]
		return p.events[0]
	} else if len(p.events) == 1 {
		p.events = make([]Event, 0)
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
	ev := last.(*OutOfStock)
	return ev.Sku() == p.SKU
}

func NewProduct(SKU string, batches []*Batch) *Product {
	return &Product{SKU: SKU, Batches: batches, VersionNumber: 0, events: make([]Event, 0)}
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
	p.events = append(p.events, NewOutOfStockEvent(line.SKU))
	return nil, OutOfStockErr
}
