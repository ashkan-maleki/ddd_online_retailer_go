package services

import (
	"context"
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
			Handle(context.Background(), event, uow.Product)
		}
	}
}

func UoW() *UnitOfWork {
	repo, err := adapters.NewProductRepo()
	log.Error(err)
	Register()
	return &UnitOfWork{Product: repo}
}
