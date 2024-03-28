package services

import (
	"context"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/link/adapters"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func newBatchRepo() *adapters.ProductRepo {
	repo, err := adapters.NewBatchRepo()
	if err != nil {
		log.Println(err)
		panic("new batch repo failed")
	}
	return repo
}

func TestAddBatch(t *testing.T) {
	ctx := context.Background()
	repo := newBatchRepo()
	ref := "b1"
	service := NewBatchService(repo)
	service.AddBatch(ctx, ref, "CRUNCHY-ARMCHAIR", 100, time.Time{})
	got := repo.Get(ctx, ref)
	assert.NotNil(t, got)
	assert.Equal(t, ref, got.Reference)
}

func TestAllocateReturnAllocation(t *testing.T) {
	ctx := context.Background()
	repo := newBatchRepo()
	ref := "b1"
	service := NewBatchService(repo)
	service.AddBatch(ctx, ref, "COMPLICATED-LAMP", 100, time.Time{})
	got, err := service.Allocate(ctx, ref, "COMPLICATED-LAMP", 10)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, ref, got)
}

func TestAllocateErrorForInvalidSku(t *testing.T) {
	ctx := context.Background()
	repo := newBatchRepo()
	ref := "b1"
	service := NewBatchService(repo)
	service.AddBatch(ctx, ref, "AREALSKU", 100, time.Time{})
	_, err := service.Allocate(ctx, ref, "NONEXISTENTSKU", 10)
	if err != nil {
		assert.ErrorIs(t, err, InvalidSku)
	}
	//assert.Equal(t, ref, got)
}
