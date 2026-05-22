package bizcache_test

import (
	"context"
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

	product, ok := bizcache.GetProductDetail(context.Background(), 1)
	if ok {
		t.Fatalf("expected cache miss when redis is nil, got hit: %+v", product)
	}

	bizcache.SetProductDetail(context.Background(), &model.Product{ID: 1, Name: "test"})
	bizcache.DeleteProductDetailCache(context.Background(), 1)
}
func TestRedisCache_IsTheKeyCorrect(t *testing.T) {
	setupTestRedis(t)
	ctx := context.Background()

	product := &model.Product{
		ID:          int64(1),
		Name:        "product detail test is the key correct",
		Description: "desc",
		PriceFen:    int64(10),
		Status:      model.ProductStatusOnSale,
	}

	bizcache.SetProductDetail(ctx, product)

	p, ok := bizcache.GetProductDetail(context.Background(), product.ID)
	if !ok {
		t.Fatalf("product cache no found")
	}

	if p.ID != product.ID && p.Name != product.Name && p.Description != product.Description && p.PriceFen != product.PriceFen && p.Status != product.Status {
		t.Fatalf("product info failed")
	}
}

func TestProductDetailCache_SetAndGetProductDetail(t *testing.T) {
	setupTestRedis(t)

	bizcache.SetProductDetail(context.Background(), &model.Product{
		ID:          1,
		Name:        "product detail on redis to set and get",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	})
	_, ok := bizcache.GetProductDetail(context.Background(), 1)
	if !ok {
		t.Fatalf("expected product detail cache exist")
	}
}

func TestProductDetailDeleteCache_GetMiss(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()

	bizcache.SetProductDetail(ctx, &model.Product{
		ID:          1,
		Name:        "product detail get miss",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	})

	bizcache.DeleteProductDetailCache(ctx, 1)
	p, ok := bizcache.GetProductDetail(ctx, 1)

	if ok {
		t.Fatalf("expected product detail cache not found")
	}

	if p != nil {
		t.Fatalf("expected product detail cache not found,got %v", p)
	}
}

func TestProductDetailCacheTTL_DoesIsTheExist(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()
	product := &model.Product{
		ID:          1,
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
}

func TestProductDetailCache_CacheHit(t *testing.T) {
	setupTestRedis(t)

	ctx := context.Background()

	_, ok := bizcache.GetProductDetail(ctx, 1)
	if ok {
		t.Fatalf("expected product detail no found")
	}

	product := &model.Product{
		ID:          1,
		Name:        "product detail on redis to set and get",
		Description: "desc",
		PriceFen:    10,
		Status:      model.ProductStatusOnSale,
	}

	bizcache.SetProductDetail(context.Background(), product)

	_, ok = bizcache.GetProductDetail(ctx, product.ID)

	if !ok {
		t.Fatalf("expected product detail exist")
	}
}

func TestProductDetail_OnSaleDeleteCache_ReturnsNil(t *testing.T) {
	setupTestRedis(t)

	// ctx := context.Background()
	// // db := global.DB

	// product := &model.Product{
	// 	ID:          int64(1),
	// 	Name:        "delete cache to on sale and off sale products",
	// 	Description: "desc",
	// 	PriceFen:    int64(100),
	// 	Status:      model.ProductStatusOffSale,
	// }

	// if err := db.Create(&product).Error; err != nil {
	// 	t.Fatalf("create product failed: %v", err)
	// }

	// if _, err := service.GetProductByID(ctx, product.ID); err != nil {
	// 	t.Fatalf("query product not found: %v", err)
	// }

	// firstProductSelect, ok := bizcache.GetProductDetail(ctx, product.ID)
	// if !ok {
	// 	t.Fatalf("expected product detail exist")
	// }

	// productOnSaleErr := service.OnSaleProduct(ctx, firstProductSelect.ID)
	// if productOnSaleErr != nil {
	// 	t.Fatalf("expected product on sale success,got %v", productOnSaleErr)
	// }

	// _, ok = bizcache.GetProductDetail(ctx, product.ID)
	// if ok {
	// 	t.Fatalf("expected product detail exist")
	// }
}

func TestProductDetail_OffSaleDeleteCache_ReturnsNil(t *testing.T) {
	setupTestRedis(t)

	// ctx := context.Background()
	db := global.DB

	product := &model.Product{
		ID:          int64(1),
		Name:        "delete cache to on sale and off sale products",
		Description: "desc",
		PriceFen:    int64(100),
		Status:      model.ProductStatusOffSale,
	}

	if err := db.Create(&product).Error; err != nil {
		t.Fatalf("create product failed: %v", err)
	}

	// if p, err := service.GetProductByID(ctx, product.ID); err != nil {
	// 	t.Fatalf("query product not found: %v", err)
	// }

	// firstProductSelect, ok := bizcache.GetProductDetail(ctx, product.ID)
	// if !ok {
	// 	t.Fatalf("expected product detail exist")
	// }

}
