package commands

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"reflect"
)

type ChangeBatchQuantity struct {
	ref string
	qty int
}

func NewChangeBatchQuantity(ref string, qty int) *ChangeBatchQuantity {
	return &ChangeBatchQuantity{ref: ref, qty: qty}
}

var _ domain.Command = (*ChangeBatchQuantity)(nil)

func (b *ChangeBatchQuantity) Qty() int {
	return b.qty
}

func (b *ChangeBatchQuantity) Ref() string {
	return b.ref
}

func (b *ChangeBatchQuantity) TransactionID() string {
	return b.ref
}

func (b *ChangeBatchQuantity) Name() string {
	eventName := reflect.TypeOf(*b).Name()
	if eventName != ChangeBatchQuantityCommand {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, ChangeBatchQuantityCommand))
	}
	return eventName
}
