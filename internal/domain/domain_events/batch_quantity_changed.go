package domain_events

import (
	"fmt"
	"reflect"
)

type BatchQuantityChanged struct {
	ref string
	qty int
}

func NewBatchQuantityChanged(ref string, qty int) *BatchQuantityChanged {
	return &BatchQuantityChanged{ref: ref, qty: qty}
}

var _ Event = (*BatchQuantityChanged)(nil)

func (b *BatchQuantityChanged) Qty() int {
	return b.qty
}

func (b *BatchQuantityChanged) Ref() string {
	return b.ref
}

func (b *BatchQuantityChanged) TransactionID() string {
	return b.ref
}

func (b *BatchQuantityChanged) Name() string {
	eventName := reflect.TypeOf(*b).Name()
	if eventName != BatchQuantityChangedEvent {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, BatchQuantityChangedEvent))
	}
	return eventName
}
