package events

import (
	"fmt"
	"reflect"
	"time"
)

type BatchCreated struct {
	sku, ref string
	qty      int
	eta      time.Time
}

func NewBatchCreated(sku string, ref string, qty int, eta time.Time) *BatchCreated {
	return &BatchCreated{sku: sku, ref: ref, qty: qty, eta: eta}
}

var _ Event = (*BatchCreated)(nil)

func (b BatchCreated) ID() string {
	return b.sku
}

func (b BatchCreated) Name() string {
	eventName := reflect.TypeOf(b).Name()
	if eventName != BatchCreatedEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, BatchCreatedEvent))
	}
	return eventName
}

func (b BatchCreated) Eta() time.Time {
	return b.eta
}

func (b BatchCreated) Qty() int {
	return b.qty
}

func (b BatchCreated) Ref() string {
	return b.ref
}

func (b BatchCreated) Sku() string {
	return b.sku
}
