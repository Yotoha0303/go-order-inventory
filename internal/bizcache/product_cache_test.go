package bizcache_test

import (
	"go-order-inventory/global"
	"go-order-inventory/internal/bizcache"
	"go-order-inventory/internal/model"
	"testing"
)

func TestProductDetailCacheKey(t *testing.T) {
	got := bizcache.ProductDetailCacheKey(123)
	want := "product:detail:123"
	if got != want {
		t.Errorf("productDetailCacheKey(123) = %q, want %q", got, want)
	}
}

func TestProductDetailCache_NoRedis(t *testing.T) {
	global.Redis = nil

	product, ok := bizcache.GetProductDetail(1)
	if ok {
		t.Fatalf("expected cache miss when redis is nil, got hit: %+v", product)
	}

	bizcache.SetProductDetail(&model.Product{ID: 1, Name: "test"})
	bizcache.DeleteProductDetailCache(1)
}
