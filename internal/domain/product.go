package domain

import "sort"

type Product struct {
	SKU           string
	VersionNumber int
	Batches       []*Batch
	Events        []Event
}

func (p *Product) HasOutOfStockEventAsLast() bool {
	if len(p.Events) == 0 {
		return false
	}
	last := p.Events[len(p.Events)-1]
	ev := last.(OutOfStockEvent)
	return ev.Sku() == p.SKU
}

func NewProduct(SKU string, batches []*Batch) *Product {
	return &Product{SKU: SKU, Batches: batches, VersionNumber: 0, Events: make([]Event, 0)}
}

func (p *Product) Allocate(line OrderLine) *Batch {
	sort.Slice(p.Batches, func(i, j int) bool {
		return p.Batches[i].ETA.Before(p.Batches[j].ETA)
	})
	for _, batch := range p.Batches {
		if batch.CanAllocate(line) {
			batch.Allocate(line)
			return batch
		}
	}
	p.Events = append(p.Events, NewOutOfStockEvent(line.SKU))
	return nil
}
