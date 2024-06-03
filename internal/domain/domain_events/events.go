package domain_events

const (
	OutOfStockEvent           = "OutOfStock"
	AllocationRequiredEvent   = "AllocationRequired"
	BatchCreatedEvent         = "BatchCreated"
	BatchQuantityChangedEvent = "BatchQuantityChanged"
)

type Event interface {
	TransactionID() string
	Name() string
}
