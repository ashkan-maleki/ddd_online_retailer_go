package services

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/gofiber/fiber/v2/log"
)

type UnitOfWork struct {
	Product *adapters.ProductRepo
}

func (uow *UnitOfWork) Commit() {
	uow.publishEvents()
}

func (uow *UnitOfWork) publishEvents() {
	for _, product := range uow.Product.Seen() {
		for product.HasEvent() {
			event := product.PopEvent()
			Handle(event)
		}
	}
}

func NewUnitOfWork() *UnitOfWork {
	repo, err := adapters.NewProductRepo()
	log.Error(err)
	return &UnitOfWork{Product: repo}
}
