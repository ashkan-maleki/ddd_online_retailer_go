package events

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
