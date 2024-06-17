package domain

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
