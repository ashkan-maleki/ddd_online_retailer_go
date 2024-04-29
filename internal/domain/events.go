package domain

import (
	"fmt"
	"reflect"
)

const (
	OutOfStockEvent         = "OutOfStock"
	AllocationRequiredEvent = "OutOfStock"
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

type AllocationRequired struct {
	orderId, sku string
	qty          int
}

func NewAllocationRequired(orderId string, sku string, qty int) *AllocationRequired {
	return &AllocationRequired{orderId: orderId, sku: sku, qty: qty}
}

var _ Event = (*AllocationRequired)(nil)

func (a AllocationRequired) Qty() int {
	return a.qty
}

func (a AllocationRequired) Sku() string {
	return a.sku
}

func (a AllocationRequired) OrderId() string {
	return a.orderId
}

func (a AllocationRequired) ID() string {
	return a.sku
}

func (a AllocationRequired) Name() string {
	eventName := reflect.TypeOf(a).Name()
	if eventName != AllocationRequiredEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, AllocationRequiredEvent))
	}
	return eventName
}
