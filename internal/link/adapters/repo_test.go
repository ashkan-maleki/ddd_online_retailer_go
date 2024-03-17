package adapters

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/domain"
	"testing"
	"time"
)

func TestRepositoryCanSaveABatch(t *testing.T) {
	batch := domain.NewBatch("batch1", "RUSTY-SOAPDISH", 100, time.Time{})
	repo, err := NewBatchRepo()
	if err != nil {
		panic(err)
	}
	repo.Add(context.Background(), *batch)
}
