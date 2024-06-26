package entity

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain/domain"
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
	//TransactionID          int64 `gorm:"primaryKey"`
	OrderLineID int64     `gorm:"primaryKey;autoIncrement:false"`
	OrderLine   OrderLine `gorm:"foreignKey:TransactionID;references:OrderLineID"`
	BatchID     int64     `gorm:"primaryKey;autoIncrement:false"`
}

type Product struct {
	SKU           string  `gorm:"primaryKey"`
	VersionNumber int     `gorm:"default:0"`
	Batches       []Batch `gorm:"foreignKey:SKU;references:SKU"`
	events        []domain.Event
}

func (p *Product) Events() []domain.Event {
	return p.events
}

func (p *Product) AddEvent(event domain.Event) {
	p.events = append(p.events, event)
}

func (p *Product) PopEvent() domain.Event {
	if len(p.events) > 1 {
		event := p.events[0]
		p.events = p.events[1:]
		return event
	} else if len(p.events) == 1 {
		event := p.events[0]
		p.events = make([]domain.Event, 0)
		return event
	} else {
		return nil
	}
}

func (p *Product) HasEvent() bool {
	return len(p.events) > 0
}
