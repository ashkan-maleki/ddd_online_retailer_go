package domain

import (
	"github.com/rs/xid"
)

type Event interface {
	ID() xid.ID
}

type OutOfStockEvent struct {
	id  xid.ID
	sku string
}

func NewOutOfStockEvent(sku string) *OutOfStockEvent {
	return &OutOfStockEvent{
		id:  xid.New(),
		sku: sku,
	}
}

var _ Event = (*OutOfStockEvent)(nil)

func (o OutOfStockEvent) ID() xid.ID {
	return o.id
}

func (o OutOfStockEvent) Sku() string {
	return o.sku
}
