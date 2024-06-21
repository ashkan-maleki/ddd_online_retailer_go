package domain

type Entity interface {
	DomainEvents() []Event
	AddDomainEvent(Event)
	HasDomainEvent() bool
	DomainEventsLength() int
}

type BaseEntity struct {
	domainEvents []Event
}

var _ Entity = (*BaseEntity)(nil)

func NewBaseEntity() *BaseEntity {
	return &BaseEntity{domainEvents: make([]Event, 0)}
}

func (p *BaseEntity) DomainEvents() []Event {
	return p.domainEvents
}

func (p *BaseEntity) AddDomainEvent(event Event) {
	p.domainEvents = append(p.domainEvents, event)
}

func (p *BaseEntity) HasDomainEvent() bool {
	return len(p.domainEvents) > 0
}

func (p *BaseEntity) DomainEventsLength() int {
	return len(p.domainEvents)
}

func (p *BaseEntity) LastEvent() Event {
	if len(p.domainEvents) == 0 {
		return nil
	}
	return p.domainEvents[len(p.domainEvents)-1]
}

func (p *BaseEntity) PopEvent() Event {
	if len(p.domainEvents) > 1 {
		p.domainEvents = p.domainEvents[1:]
		return p.domainEvents[0]
	} else if len(p.domainEvents) == 1 {
		p.domainEvents = make([]Event, 0)
		return p.domainEvents[0]
	} else {
		return nil
	}
}
