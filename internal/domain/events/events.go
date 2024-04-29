package events

const (
	OutOfStockEvent           = "OutOfStock"
	AllocationRequiredEvent   = "AllocationRequired"
	BatchCreatedEvent         = "BatchCreated"
	BatchQuantityChangedEvent = "BatchQuantityChanged"
)

type Event interface {
	ID() string
	Name() string
}
