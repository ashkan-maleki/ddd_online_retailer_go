package domain_events

import (
	"fmt"
	"reflect"
	"time"
)

type BatchCreated struct {
	Sku, ref string
	qty      int
	eta      time.Time
}

func NewBatchCreated(sku string, ref string, qty int, eta time.Time) *BatchCreated {
	return &BatchCreated{Sku: sku, ref: ref, qty: qty, eta: eta}
}

var _ Event = (*BatchCreated)(nil)

func (e *BatchCreated) TransactionID() string {
	return e.Sku
}

func (e *BatchCreated) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != BatchCreatedEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, BatchCreatedEvent))
	}
	return eventName
}

func (e *BatchCreated) Eta() time.Time {
	return e.eta
}

func (e *BatchCreated) Qty() int {
	return e.qty
}

func (e *BatchCreated) Ref() string {
	return e.ref
}

func (e *BatchCreated) GetSku() string {
	return e.Sku
}
