package main

import (
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/presentation"
)

func main() {
	app, err := presentation.NewApplication()
	if err != nil {
		panic(err)
	}
	app.Start()
}
