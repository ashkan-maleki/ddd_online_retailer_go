package presentation

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/services"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Application struct {
	service *services.BatchService
	rest    *fiber.App
}

func NewApplication() (*Application, error) {

	app := &Application{}
	app.rest = fiber.New()
	repo, err := adapters.NewBatchRepo()
	if err != nil {
		return nil, err
	}
	app.service = services.NewBatchService(repo)
	return app, nil
}

func (app Application) Start() {

	app.rest.Post("/add-batch", func(c *fiber.Ctx) error {
		batch := new(Batch)
		err := c.BodyParser(batch)
		if err != nil {
			return err
		}
		app.service.AddBatch(c.Context(), batch.Reference, batch.SKU, batch.PurchasedQuantity, batch.ETA)
		return c.Status(201).SendString("batch is added")
	})
	app.rest.Get("/allocate", func(c *fiber.Ctx) error {
		line := new(OrderLine)
		err := c.BodyParser(line)
		if err != nil {
			return err
		}
		allocate, err := app.service.Allocate(c.Context(), line.OrderID, line.SKU, line.Qty)
		if err != nil {
			return err
		}
		return c.Status(201).SendString(fmt.Sprintf("line is allocated: %v", allocate))
	})

	log.Println("http://127.0.0.1:3000")
	log.Fatal(app.rest.Listen(":3000"))
}
