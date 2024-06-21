package domain_events

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"reflect"
)

type Allocated struct {
	orderID, sku, batchRef string
	qty                    int
}

func NewAllocated(orderId, sku, batchRef string, qty int) *Allocated {
	return &Allocated{orderID: orderId, sku: sku, batchRef: batchRef, qty: qty}
}

var _ domain.Event = (*Allocated)(nil)

func (e *Allocated) Qty() int {
	return e.qty
}

func (e *Allocated) SKU() string {
	return e.sku
}

func (e *Allocated) OrderID() string {
	return e.orderID
}

func (e *Allocated) BatchRef() string {
	return e.batchRef
}

func (e *Allocated) TransactionID() string {
	return e.sku
}

func (e *Allocated) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != AllocatedEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, AllocatedEvent))
	}
	return eventName
}

func ConvertAllocate(event domain.Event) (*Allocated, error) {
	//var allocated *Allocated
	//switch a := event.(type) {
	//case *Allocated:
	//	allocated = a
	//	break
	//default:
	//	return nil, fmt.Errorf("wrong event type %v", event.Name())
	//}
	//return allocated, nil
	return domain.Convert[*Allocated](event)
}
