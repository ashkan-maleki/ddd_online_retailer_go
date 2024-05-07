package services

//func newBatchRepo() *adapters.ProductRepo {
//	repo, err := adapters.NewProductRepo()
//	if err != nil {
//		log.Println(err)
//		panic("new batch repo failed")
//	}
//	return repo
//}
//
//func TestAddBatch(t *testing.T) {
//	ctx := context.Background()
//	repo := newBatchRepo()
//	service := NewBatchService(repo)
//
//	ref := "b1"
//	sku := "CRUNCHY-ARMCHAIR"
//
//	_ = service.AddBatch(ctx, ref, sku, 100, time.Time{})
//	product := repo.Get(ctx, sku)
//	assert.NotNil(t, product)
//	assert.NotNil(t, product.Batches)
//	found := false
//	for _, batch := range product.Batches {
//		if batch.Reference == ref {
//			found = true
//		}
//	}
//	assert.True(t, found)
//}

//func TestAllocateReturnAllocation(t *testing.T) {
//	ctx := context.Background()
//	repo := newBatchRepo()
//	ref := "b1"
//	service := NewBatchService(repo)
//	service.AddBatch(ctx, ref, "COMPLICATED-LAMP", 100, time.Time{})
//	got, err := service.Allocate(ctx, ref, "COMPLICATED-LAMP", 10)
//	if err != nil {
//		panic(err)
//	}
//	assert.Equal(t, ref, got)
//}

//func TestAllocateErrorForInvalidSku(t *testing.T) {
//	ctx := context.Background()
//	repo := newBatchRepo()
//	ref := "b1"
//	service := NewBatchService(repo)
//	service.AddBatch(ctx, ref, "AREALSKU", 100, time.Time{})
//	_, err := service.Allocate(ctx, ref, "NONEXISTENTSKU", 10)
//	if err != nil {
//		assert.ErrorIs(t, err, InvalidSku)
//	}
//	//assert.Equal(t, ref, got)
//}
