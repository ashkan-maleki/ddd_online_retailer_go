package commands

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"reflect"
)

type Allocate struct {
	orderId, sku string
	qty          int
}

func NewAllocated(orderId string, sku string, qty int) *Allocate {
	return &Allocate{orderId: orderId, sku: sku, qty: qty}
}

var _ domain.Command = (*Allocate)(nil)

func (e *Allocate) Qty() int {
	return e.qty
}

func (e *Allocate) Sku() string {
	return e.sku
}

func (e *Allocate) OrderId() string {
	return e.orderId
}

func (e *Allocate) TransactionID() string {
	return e.sku
}

func (e *Allocate) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != AllocateCommand {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, AllocateCommand))
	}
	return eventName
}
