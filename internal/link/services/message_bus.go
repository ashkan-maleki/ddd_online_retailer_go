package services

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
)

type HandleFunc func(event events.Event)

var Handlers = make(map[string][]HandleFunc)

func SendOutOfStockNotification(event events.Event) {
	outOfStock, ok := event.(events.OutOfStock)
	if ok {
		panic(fmt.Sprintf("wrong event type %v", event.Name()))
	}
	SendEmail("stock@eshop.com", fmt.Sprintf("out of stock for %v", outOfStock.Sku()))
}

func AddBatch(event events.Event) {

}

func Allocate(event events.Event) {

}

func ChangeBatchQuantity(event events.Event) {

}

func Register() {
	Handlers[events.BatchCreatedEvent] = []HandleFunc{AddBatch}
	Handlers[events.BatchQuantityChangedEvent] = []HandleFunc{ChangeBatchQuantity}
	Handlers[events.AllocationRequiredEvent] = []HandleFunc{Allocate}
	Handlers[events.OutOfStockEvent] = []HandleFunc{SendOutOfStockNotification}
}

func Handle(event events.Event) {
	handlers, ok := Handlers[event.Name()]
	if !ok {
		panic(fmt.Sprintf("no handler is registered for %v", event.Name()))
	}
	for _, handler := range handlers {
		handler(event)
	}
}
