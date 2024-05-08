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
	eves := make([]events.Event, 0)
	for _, product := range repo.Seen() {
		for product.HasEvent() {
			event := product.PopEvent()
			eves = append(eves, event)
		}
	}
	fmt.Println("collected size, inside: ", len(eves))
	return eves
}

type iterator func(yield func(events.Event) bool)

func Handle(ctx context.Context, event events.Event, repo *adapters.ProductRepo) ([]any, error) {

	results := make([]any, 0)
	queue := []events.Event{event}
	handlersErrors := make([]error, 0)
	for len(queue) > 0 {
		eventInQueue := queue[0]
		fmt.Println("Handle, event ouf queue: ", fmt.Sprintf("(%v, %T)", eventInQueue, eventInQueue))
		queue = queue[1:]
		fmt.Println("event name: ", eventInQueue.Name())
		handlers, ok := Handlers[eventInQueue.Name()]
		if !ok {
			return nil, fmt.Errorf("no handler is registered for %v", eventInQueue.Name())
		}
		for _, handler := range handlers {
			result, err := handler(ctx, eventInQueue, repo)
			if err != nil {
				handlersErrors = append(handlersErrors, err)
				//return nil, err
			}
			if result != nil {
				results = append(results, result)
			}
			//fmt.Println("collected size, outside: ", len(collectNewEvents(repo)))
			collectedEvents := collectNewEvents(repo)
			for _, ev := range collectedEvents {
				//fmt.Println("Handle, events: ", fmt.Sprintf("(%v, %T)", ev, ev))
				queue = append(queue, ev)
			}
		}
		fmt.Println("queue size: ", len(queue))

	}
	var err error
	if len(handlersErrors) > 0 {
		err = handlersErrors[0]
	}
	//fmt.Println("event loop: ", fmt.Sprintf("(%v, %T)", event, event))
	return results, err
}
