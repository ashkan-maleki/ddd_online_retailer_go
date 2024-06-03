package domain_events

import (
	"fmt"
	"reflect"
)

type AllocationRequired struct {
	orderId, sku string
	qty          int
}

func NewAllocationRequired(orderId string, sku string, qty int) *AllocationRequired {
	return &AllocationRequired{orderId: orderId, sku: sku, qty: qty}
}

var _ Event = (*AllocationRequired)(nil)

func (e *AllocationRequired) Qty() int {
	return e.qty
}

func (e *AllocationRequired) Sku() string {
	return e.sku
}

func (e *AllocationRequired) OrderId() string {
	return e.orderId
}

func (e *AllocationRequired) TransactionID() string {
	return e.sku
}

func (e *AllocationRequired) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != AllocationRequiredEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, AllocationRequiredEvent))
	}
	return eventName
}
