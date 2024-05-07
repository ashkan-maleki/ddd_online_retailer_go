package services

import (
	"context"
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
)

type HandleFunc func(context.Context, events.Event, *adapters.ProductRepo) (any, error)

var Handlers = make(map[string][]HandleFunc)

func Register() {
	Handlers[events.BatchCreatedEvent] = []HandleFunc{AddBatch}
	Handlers[events.BatchQuantityChangedEvent] = []HandleFunc{ChangeBatchQuantity}
	Handlers[events.AllocationRequiredEvent] = []HandleFunc{Allocate}
	Handlers[events.OutOfStockEvent] = []HandleFunc{SendOutOfStockNotification}
}

func collectNewEventsChannel(repo *adapters.ProductRepo) <-chan events.Event {
	ch := make(chan events.Event)
	go func() {
		for _, product := range repo.Seen() {
			for product.HasEvent() {
				event := product.PopEvent()
				ch <- event
			}
		}
	}()
	return ch
}

func collectNewEvents(repo *adapters.ProductRepo) []events.Event {
	fmt.Println("***********")
	fmt.Println("seen: ", repo.Seen())
	eves := make([]events.Event, 0)
	for _, product := range repo.Seen() {
		fmt.Println("product: ", product.SKU)
		fmt.Println("product event size: ", len(product.Events()))
		for product.HasEvent() {
			event := product.PopEvent()
			fmt.Println("event loop: ", event)
			eves = append(eves, event)
		}
	}
	fmt.Println("collected events size: ", len(eves))
	fmt.Println("collected events: ", eves)
	return eves
}

type iterator func(yield func(events.Event) bool)

func Handle(ctx context.Context, event events.Event, repo *adapters.ProductRepo) ([]any, error) {
	fmt.Println("event name: ", event.Name())
	handlers, ok := Handlers[event.Name()]
	if !ok {
		return nil, fmt.Errorf("no handler is registered for %v", event.Name())
	}
	results := make([]any, 0)
	queue := []events.Event{event}
	handlersErrors := make([]error, 0)
	for len(queue) > 0 {
		eventInQueue := queue[0]
		for _, handler := range handlers {
			result, err := handler(ctx, eventInQueue, repo)
			if err != nil {
				handlersErrors = append(handlersErrors, err)
				//return nil, err
			}
			results = append(results, result)
			for _, ev := range collectNewEvents(repo) {
				fmt.Println("events: ", ev)
				queue = append(queue, ev)
			}
		}
		queue = queue[1:]
	}
	var err error
	if len(handlersErrors) > 0 {
		err = handlersErrors[0]
	}
	return results, err
}
