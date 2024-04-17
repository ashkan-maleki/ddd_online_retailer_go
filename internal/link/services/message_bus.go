package services

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
)

type HandleFunc func(event domain.Event)

var Handlers = make(map[string][]HandleFunc)

func SendOutOfStockNotification(event domain.Event) {
	outOfStock, ok := event.(domain.OutOfStock)
	if ok {
		panic(fmt.Sprintf("wrong event type %v", event.Name()))
	}
	SendEmail("stock@eshop.com", fmt.Sprintf("out of stock for %v", outOfStock.Sku()))
}

func Register() {
	Handlers[domain.OutOfStockEvent] = []HandleFunc{SendOutOfStockNotification}
}

func Handle(event domain.Event) {
	handlers, ok := Handlers[event.Name()]
	if !ok {
		panic(fmt.Sprintf("no handler is registered for %v", event.Name()))
	}
	for _, handler := range handlers {
		handler(event)
	}
}

// TODO: Read