package adapters

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCanConvertDbBatchToDomainBatch(t *testing.T) {
	now := time.Now()
	db := Batches{
		Id:        1,
		Sku:       "my-sku",
		Reference: "my-ref",
		Eta:       now,
	}
	dom := ToDomainBatch(db)
	assert.Equal(t, db.Reference, dom.Reference)
	assert.Equal(t, db.Sku, dom.SKU)
	assert.Equal(t, db.Eta, dom.ETA)
}
