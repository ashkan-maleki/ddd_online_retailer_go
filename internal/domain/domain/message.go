package domain

import "fmt"

type Message interface {
	TransactionID() string
	Name() string
}

type Command interface {
	Message
}

type Event interface {
	Message
}

func Convert[E Event](event Event) (E, error) {
	switch a := event.(type) {
	case E:
		return a, nil
	default:
		var eve E
		return eve, fmt.Errorf("wrong event type %v", event.Name())
	}

}
