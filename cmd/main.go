package main

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"time"
)

func main() {
	fmt.Println("welcome to online retailer")
	batch := domain.NewBatch("batch-001", "SMALL-TABLE", 20, time.Now())
	line := domain.NewOrderLine("order-ref", "SMALL-TABLE", 2)
	batch.Allocate(line)
}
