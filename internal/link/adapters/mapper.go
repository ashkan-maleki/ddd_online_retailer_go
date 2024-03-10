package adapters

import (
	"fmt"
	"github.com/ashkan-maleki/ddd_online_retailer_go/internal/ddd"
	"reflect"
)

func ToDomainBatch(dbBatch Batches) ddd.Batch {
	var batch ddd.Batch
	return batch
}

func a() {
	type T struct {
		A int
		B string
	}
	t := T{23, "skidoo"}
	s := reflect.ValueOf(&t).Elem().FieldByNameFunc(func(s string) bool {
		return true
	})
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
