package domain

import (
	"fmt"
	"reflect"
)

const (
	OutOfStockEvent = "OutOfStock"
)

type Event interface {
	ID() string
	Name() string
}

type OutOfStock struct {
	sku string
}

func NewOutOfStockEvent(sku string) *OutOfStock {
	return &OutOfStock{
		sku: sku,
	}

}

var _ Event = (*OutOfStock)(nil)

func (o OutOfStock) ID() string {
	return o.sku
}

func (o OutOfStock) Sku() string {
	return o.sku
}

func (o OutOfStock) Name() string {

	eventName := reflect.TypeOf(o).Name()
	if eventName != OutOfStockEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, OutOfStockEvent))
	}
	return eventName
}
