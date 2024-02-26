package main

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/ddd"
	"time"
)

func main() {
	fmt.Println("welcome to online retailer")
	batch := ddd.NewBatch("batch-001", "SMALL-TABLE", time.Now(), 20)
	line := ddd.NewOrderLine("order-ref", "SMALL-TABLE", 2)
	batch.Allocate(line)
}
