package domain_events

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"reflect"
)

type Deallocated struct {
	orderId, sku string
	qty          int
}

func NewDeallocated(orderId, sku string, qty int) *Deallocated {
	return &Deallocated{orderId: orderId, sku: sku, qty: qty}
}

var _ domain.Command = (*Deallocated)(nil)

func (e *Deallocated) Qty() int {
	return e.qty
}

func (e *Deallocated) Sku() string {
	return e.sku
}

func (e *Deallocated) OrderId() string {
	return e.orderId
}

func (e *Deallocated) TransactionID() string {
	return e.sku
}

func (e *Deallocated) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != DeallocatedEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, DeallocatedEvent))
	}
	return eventName
}
