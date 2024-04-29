package entity

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/events"
	"time"
)

type OrderLine struct {
	ID      int64 `gorm:"primaryKey"`
	OrderID string
	SKU     string
	Qty     int
	//Allocations []Allocation `gorm:"foreignKey:OrderLineID"`
}

type Batch struct {
	ID                int64 `gorm:"primaryKey"`
	Reference         string
	SKU               string
	Product           Product `gorm:"foreignKey:SKU;references:SKU"`
	PurchasedQuantity int
	ETA               time.Time
	Allocations       []Allocation `gorm:"foreignKey:BatchID"`
}

type Allocation struct {
	ID          int64 `gorm:"primaryKey"`
	OrderLineID int64
	OrderLine   OrderLine `gorm:"foreignKey:ID;references:OrderLineID"`
	BatchID     int64
}

type Product struct {
	SKU           string  `gorm:"primaryKey"`
	VersionNumber int     `gorm:"default:0"`
	Batches       []Batch `gorm:"foreignKey:SKU;references:SKU"`
	events        []events.Event
}

func (p *Product) Events() []events.Event {
	return p.events
}

func (p *Product) AddEvent(event events.Event) {
	p.events = append(p.events, event)
}

func (p *Product) PopEvent() events.Event {
	if len(p.events) > 1 {
		p.events = p.events[1:]
		return p.events[0]
	} else if len(p.events) == 1 {
		p.events = make([]events.Event, 0)
		return p.events[0]
	} else {
		return nil
	}
}

func (p *Product) HasEvent() bool {
	return len(p.events) > 0
}
