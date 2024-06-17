package commands

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
	"reflect"
	"time"
)

type CreateBatch struct {
	Sku, ref string
	qty      int
	eta      time.Time
}

func NewCreateBatch(sku string, ref string, qty int, eta time.Time) *CreateBatch {
	return &CreateBatch{Sku: sku, ref: ref, qty: qty, eta: eta}
}

var _ domain.Command = (*CreateBatch)(nil)

func (e *CreateBatch) TransactionID() string {
	return e.Sku
}

func (e *CreateBatch) Name() string {
	eventName := reflect.TypeOf(*e).Name()
	if eventName != CreateBatchCommand {
		panic(fmt.Sprintf("event name collision, %v != %v", eventName, CreateBatchCommand))
	}
	return eventName
}

func (e *CreateBatch) Eta() time.Time {
	return e.eta
}

func (e *CreateBatch) Qty() int {
	return e.qty
}

func (e *CreateBatch) Ref() string {
	return e.ref
}

func (e *CreateBatch) GetSku() string {
	return e.Sku
}
