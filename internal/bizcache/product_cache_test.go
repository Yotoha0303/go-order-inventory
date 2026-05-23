package bizcache_test

import (
	"context"
	"go-order-inventory/global"
	"go-order-inventory/internal/bizcache"
	"go-order-inventory/internal/model"
	"testing"
)

func TestProductDetailCacheKey(t *testing.T) {

	got := bizcache.ProductDetailCacheKey(1001)
	want := "product:detail:1001"
	if got != want {
		t.Errorf("productDetailCacheKey(1001) = %q, want %q", got, want)
	}

}

func TestProductDetailCache_NoRedis(t *testing.T) {
	oldRedis := global.Redis
	global.Redis = nil
	defer func() {
		global.Redis = oldRedis
	}()

	product, ok := bizcache.GetProductDetail(context.Background(), 1002)
	if ok {
		t.Fatalf("expected cache miss when redis is nil, got hit: %+v", product)
	}

	bizcache.SetProductDetail(context.Background(), &model.Product{ID: 1002, Name: "test"})
	bizcache.DeleteProductDetailCache(context.Background(), 1002)
}

func TestProductDetailCache_SetGet_WithRedis(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	product := &model.Product{
		ID:          1004,
		Name:        "product detail on redis to set and get",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	}

	bizcache.SetProductDetail(ctx, product)
	_, ok := bizcache.GetProductDetail(ctx, product.ID)
	if !ok {
		t.Fatalf("expected product detail cache exist")
	}

	bizcache.DeleteProductDetailCache(ctx, product.ID)
}

func TestProductDetailCache_DeleteMiss_WithRedis(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	product := &model.Product{
		ID:          1005,
		Name:        "product detail get miss",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	}

	bizcache.SetProductDetail(ctx, product)
	p, ok := bizcache.GetProductDetail(ctx, product.ID)

	if !ok {
		t.Fatalf("expected product detail cache found")
	}

	if p.ID != product.ID || p.Name != product.Name || p.Description != product.Description || p.PriceFen != product.PriceFen || p.Status != product.Status {
		t.Fatalf("product info failed")
	}

	bizcache.DeleteProductDetailCache(ctx, product.ID)
	p, ok = bizcache.GetProductDetail(ctx, product.ID)

	if ok {
		t.Fatalf("expected product detail cache not found")
	}

	if p != nil {
		t.Fatalf("expected product detail cache not found,got %v", p)
	}

	bizcache.DeleteProductDetailCache(context.Background(), product.ID)
}

func TestProductDetailCache_TTL_WithRedis(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	product := &model.Product{
		ID:          1006,
		Name:        "product detail ttl",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	}

	bizcache.SetProductDetail(ctx, product)

	ttl, err := global.Redis.TTL(ctx, bizcache.ProductDetailCacheKey(product.ID)).Result()
	if err != nil {
		t.Fatalf("query ttl failed: %v", err)
	}
	if ttl <= 0 {
		t.Fatalf("expected ttl > 0, got %v", ttl)
	}
	if ttl > bizcache.ProductDetailCacheTTL {
		t.Fatalf("expected ttl <= %v, got %v", bizcache.ProductDetailCacheTTL, ttl)
	}

	bizcache.DeleteProductDetailCache(context.Background(), product.ID)
}
